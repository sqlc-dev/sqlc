//go:build !nowasm && cgo && ((linux && amd64) || (linux && arm64) || (darwin && amd64) || (darwin && arm64) || (windows && amd64))

// The above build constraint is based of the cgo directives in this file:
// https://github.com/bytecodealliance/wasmtime-go/blob/main/ffi.go
package wasm

import (
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
	"runtime/trace"
	"strings"

	wasmtime "github.com/bytecodealliance/wasmtime-go/v14"
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

// This version must be updated whenever the wasmtime-go dependency is updated
const wasmtimeVersion = `v14.0.0`

func cacheDir() (string, error) {
	cache := os.Getenv("SQLCCACHE")
	if cache != "" {
		return cache, nil
	}
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		cacheHome = filepath.Join(home, ".cache")
	}
	return filepath.Join(cacheHome, "sqlc"), nil
}

var flight singleflight.Group

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

func (r *Runner) loadModule(ctx context.Context, engine *wasmtime.Engine) (*wasmtime.Module, error) {
	expected, err := r.getChecksum(ctx)
	if err != nil {
		return nil, err
	}
	value, err, _ := flight.Do(expected, func() (interface{}, error) {
		return r.loadSerializedModule(ctx, engine, expected)
	})
	if err != nil {
		return nil, err
	}
	data, ok := value.([]byte)
	if !ok {
		return nil, fmt.Errorf("returned value was not a byte slice")
	}
	return wasmtime.NewModuleDeserialize(engine, data)
}

func (r *Runner) loadSerializedModule(ctx context.Context, engine *wasmtime.Engine, expectedSha string) ([]byte, error) {
	cacheDir, err := cache.PluginsDir()
	if err != nil {
		return nil, err
	}

	pluginDir := filepath.Join(cacheDir, expectedSha)
	modName := fmt.Sprintf("plugin_%s_%s_%s.module", runtime.GOOS, runtime.GOARCH, wasmtimeVersion)
	modPath := filepath.Join(pluginDir, modName)
	_, staterr := os.Stat(modPath)
	if staterr == nil {
		data, err := os.ReadFile(modPath)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	wmod, err := r.loadWASM(ctx, cacheDir, expectedSha)
	if err != nil {
		return nil, err
	}

	moduRegion := trace.StartRegion(ctx, "wasmtime.NewModule")
	module, err := wasmtime.NewModule(engine, wmod)
	moduRegion.End()
	if err != nil {
		return nil, fmt.Errorf("define wasi: %w", err)
	}

	err = os.Mkdir(pluginDir, 0755)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("mkdirall: %w", err)
	}
	out, err := module.Serialize()
	if err != nil {
		return nil, fmt.Errorf("serialize: %w", err)
	}
	if err := os.WriteFile(modPath, out, 0444); err != nil {
		return nil, fmt.Errorf("cache wasm: %w", err)
	}

	return out, nil
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

func (r *Runner) loadWASM(ctx context.Context, cache string, expected string) ([]byte, error) {
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

	return wmod, nil
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

	engine := wasmtime.NewEngine()
	module, err := r.loadModule(ctx, engine)
	if err != nil {
		return fmt.Errorf("loadModule: %w", err)
	}

	linker := wasmtime.NewLinker(engine)
	if err := linker.DefineWasi(); err != nil {
		return err
	}

	dir, err := os.MkdirTemp(os.Getenv("SQLCTMPDIR"), "out")
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}

	defer os.RemoveAll(dir)
	stdinPath := filepath.Join(dir, "stdin")
	stderrPath := filepath.Join(dir, "stderr")
	stdoutPath := filepath.Join(dir, "stdout")

	if err := os.WriteFile(stdinPath, stdinBlob, 0755); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	// Configure WASI imports to write stdout into a file.
	wasiConfig := wasmtime.NewWasiConfig()
	wasiConfig.SetArgv([]string{"plugin.wasm", method})
	wasiConfig.SetStdinFile(stdinPath)
	wasiConfig.SetStdoutFile(stdoutPath)
	wasiConfig.SetStderrFile(stderrPath)

	keys := []string{"SQLC_VERSION"}
	vals := []string{info.Version}
	for _, key := range r.Env {
		keys = append(keys, key)
		vals = append(vals, os.Getenv(key))
	}
	wasiConfig.SetEnv(keys, vals)

	store := wasmtime.NewStore(engine)
	store.SetWasi(wasiConfig)

	linkRegion := trace.StartRegion(ctx, "linker.DefineModule")
	err = linker.DefineModule(store, "", module)
	linkRegion.End()
	if err != nil {
		return fmt.Errorf("define wasi: %w", err)
	}

	// Run the function
	fn, err := linker.GetDefault(store, "")
	if err != nil {
		return fmt.Errorf("wasi: get default: %w", err)
	}

	callRegion := trace.StartRegion(ctx, "call _start")
	_, err = fn.Call(store)
	callRegion.End()

	if cerr := checkError(err, stderrPath); cerr != nil {
		return cerr
	}

	// Print WASM stdout
	stdoutBlob, err := os.ReadFile(stdoutPath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

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

func checkError(err error, stderrPath string) error {
	if err == nil {
		return err
	}

	var wtError *wasmtime.Error
	if errors.As(err, &wtError) {
		if code, ok := wtError.ExitStatus(); ok {
			if code == 0 {
				return nil
			}
		}
	}
	// Print WASM stdout
	stderrBlob, rferr := os.ReadFile(stderrPath)
	if rferr == nil && len(stderrBlob) > 0 {
		return errors.New(string(stderrBlob))
	}
	return fmt.Errorf("call: %w", err)
}
