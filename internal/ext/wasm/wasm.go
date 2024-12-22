package wasm

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/sqlc-dev/sqlc/internal/cache"
	"github.com/sqlc-dev/sqlc/internal/info"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

var flight singleflight.Group

type runtimeAndCode struct {
	rt   wazero.Runtime
	code wazero.CompiledModule
}

// Verify the provided sha256 is valid.
func (r *Runner) getChecksum(ctx context.Context) (string, error) {
	if r.SHA256 != "" {
		return r.SHA256, nil
	}
	// TODO: Add a log line here about something
	_, sum, err := r.fetch(ctx, r.URL)
	if err != nil {
		return "", err
	}
	slog.Warn("fetching WASM binary to calculate sha256. Set this value in sqlc.yaml to prevent unneeded work", "sha256", sum)
	return sum, nil
}

func (r *Runner) loadAndCompile(ctx context.Context) (*runtimeAndCode, error) {
	expected, err := r.getChecksum(ctx)
	if err != nil {
		return nil, err
	}
	cacheDir, err := cache.PluginsDir()
	if err != nil {
		return nil, err
	}
	value, err, _ := flight.Do(expected, func() (interface{}, error) {
		return r.loadAndCompileWASM(ctx, cacheDir, expected)
	})
	if err != nil {
		return nil, err
	}
	data, ok := value.(*runtimeAndCode)
	if !ok {
		return nil, fmt.Errorf("returned value was not a compiled module")
	}
	return data, nil
}

func (r *Runner) fetch(ctx context.Context, uri string) ([]byte, string, error) {
	var body io.ReadCloser

	switch {

	case strings.HasPrefix(uri, "file://"):
		file, err := os.Open(strings.TrimPrefix(uri, "file://"))
		if err != nil {
			return nil, "", fmt.Errorf("os.Open: %s %w", uri, err)
		}
		body = file

	case strings.HasPrefix(uri, "https://"):
		req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
		if err != nil {
			return nil, "", fmt.Errorf("http.Get: %s %w", uri, err)
		}
		req.Header.Set("User-Agent", fmt.Sprintf("sqlc/%s Go/%s (%s %s)", info.Version, runtime.Version(), runtime.GOOS, runtime.GOARCH))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, "", fmt.Errorf("http.Get: %s %w", r.URL, err)
		}
		body = resp.Body

	default:
		return nil, "", fmt.Errorf("unknown scheme: %s", r.URL)
	}

	defer body.Close()

	wmod, err := io.ReadAll(body)
	if err != nil {
		return nil, "", fmt.Errorf("readall: %w", err)
	}

	sum := sha256.Sum256(wmod)
	actual := fmt.Sprintf("%x", sum)

	return wmod, actual, nil
}

func (r *Runner) loadAndCompileWASM(ctx context.Context, cache string, expected string) (*runtimeAndCode, error) {
	pluginDir := filepath.Join(cache, expected)
	pluginPath := filepath.Join(pluginDir, "plugin.wasm")
	_, staterr := os.Stat(pluginPath)

	uri := r.URL
	if staterr == nil {
		uri = "file://" + pluginPath
	}

	wmod, actual, err := r.fetch(ctx, uri)
	if err != nil {
		return nil, err
	}

	if expected != actual {
		return nil, fmt.Errorf("invalid checksum: expected %s, got %s", expected, actual)
	}

	if staterr != nil {
		err := os.Mkdir(pluginDir, 0755)
		if err != nil && !os.IsExist(err) {
			return nil, fmt.Errorf("mkdirall: %w", err)
		}
		if err := os.WriteFile(pluginPath, wmod, 0444); err != nil {
			return nil, fmt.Errorf("cache wasm: %w", err)
		}
	}

	wazeroCache, err := wazero.NewCompilationCacheWithDir(filepath.Join(cache, "wazero"))
	if err != nil {
		return nil, fmt.Errorf("wazero.NewCompilationCacheWithDir: %w", err)
	}

	config := wazero.NewRuntimeConfig().WithCompilationCache(wazeroCache)
	rt := wazero.NewRuntimeWithConfig(ctx, config)

	if _, err := wasi_snapshot_preview1.Instantiate(ctx, rt); err != nil {
		return nil, fmt.Errorf("wasi_snapshot_preview1 instantiate: %w", err)
	}

	// Compile the Wasm binary once so that we can skip the entire compilation
	// time during instantiation.
	code, err := rt.CompileModule(ctx, wmod)
	if err != nil {
		return nil, fmt.Errorf("compile module: %w", err)
	}

	return &runtimeAndCode{rt: rt, code: code}, nil
}

// removePGCatalog removes the pg_catalog schema from the request. There is a
// mysterious (reason unknown) bug with wasm plugins when a large amount of
// tables (like there are in the catalog) are sent.
// @see https://github.com/sqlc-dev/sqlc/pull/1748
func removePGCatalog(req *plugin.GenerateRequest) {
	if req.Catalog == nil || req.Catalog.Schemas == nil {
		return
	}

	filtered := make([]*plugin.Schema, 0, len(req.Catalog.Schemas))
	for _, schema := range req.Catalog.Schemas {
		if schema.Name == "pg_catalog" || schema.Name == "information_schema" {
			continue
		}

		filtered = append(filtered, schema)
	}

	req.Catalog.Schemas = filtered
}

func (r *Runner) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	req, ok := args.(protoreflect.ProtoMessage)
	if !ok {
		return status.Error(codes.InvalidArgument, "args isn't a protoreflect.ProtoMessage")
	}

	// Remove the pg_catalog schema. Its sheer size causes unknown issues with wasm plugins
	genReq, ok := req.(*plugin.GenerateRequest)
	if ok {
		removePGCatalog(genReq)
		req = genReq
	}

	stdinBlob, err := proto.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to encode codegen request: %w", err)
	}

	runtimeAndCode, err := r.loadAndCompile(ctx)
	if err != nil {
		return fmt.Errorf("loadBytes: %w", err)
	}

	var stderr, stdout bytes.Buffer

	conf := wazero.NewModuleConfig().
		WithName("").
		WithArgs("plugin.wasm", method).
		WithStdin(bytes.NewReader(stdinBlob)).
		WithStdout(&stdout).
		WithStderr(&stderr).
		WithEnv("SQLC_VERSION", info.Version)
	for _, key := range r.Env {
		conf = conf.WithEnv(key, os.Getenv(key))
	}

	result, err := runtimeAndCode.rt.InstantiateModule(ctx, runtimeAndCode.code, conf)
	if err == nil {
		defer result.Close(ctx)
	}
	if cerr := checkError(err, stderr); cerr != nil {
		return cerr
	}

	// Print WASM stdout
	stdoutBlob := stdout.Bytes()

	resp, ok := reply.(protoreflect.ProtoMessage)
	if !ok {
		return fmt.Errorf("reply isn't a GenerateResponse")
	}

	if err := proto.Unmarshal(stdoutBlob, resp); err != nil {
		return err
	}

	return nil
}

func (r *Runner) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func checkError(err error, stderr bytes.Buffer) error {
	if err == nil {
		return err
	}

	if exitErr, ok := err.(*sys.ExitError); ok {
		if exitErr.ExitCode() == 0 {
			return nil
		}
	}

	// Print WASM stdout
	stderrBlob := stderr.String()
	if len(stderrBlob) > 0 {
		return errors.New(stderrBlob)
	}
	return fmt.Errorf("call: %w", err)
}
