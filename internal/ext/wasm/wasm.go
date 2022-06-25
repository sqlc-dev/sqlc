//go:build cgo && ((linux && amd64) || (linux && arm64) || (darwin && amd64) || (darwin && arm64) || (windows && amd64))

// The above build constraint is based of the cgo directives in this file:
// https://github.com/bytecodealliance/wasmtime-go/blob/main/ffi.go
package wasm

import (
	"context"
	"crypto/sha256"
	_ "embed"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/trace"
	"strings"

	wasmtime "github.com/bytecodealliance/wasmtime-go"

	"github.com/kyleconroy/sqlc/internal/plugin"
)

func cacheDir() (string, error) {
	// Use the checksum to see if it already existsin the modcache
	cache := os.Getenv("SQLCCACHE")
	if cache != "" {
		return cache, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "sqlc", "mod"), nil
}

type Runner struct {
	URL      string
	Checksum string
}

// Verify the provided sha256 is valid.
func (r *Runner) parseChecksum() (string, error) {
	if r.Checksum == "" {
		return "", fmt.Errorf("missing checksum")
	}
	if !strings.HasPrefix(r.Checksum, "sha256/") {
		return "", fmt.Errorf("invalid checksum algo: %s", r.Checksum)
	}
	return strings.TrimPrefix(r.Checksum, "sha256/"), nil
}

func (r *Runner) loadModule(ctx context.Context, engine *wasmtime.Engine) (*wasmtime.Module, error) {
	expected, err := r.parseChecksum()
	if err != nil {
		return nil, err
	}

	cache, err := cacheDir()
	if err != nil {
		return nil, err
	}

	pluginDir := filepath.Join(cache, expected)
	// TODO: Include os / arch in module name
	modPath := filepath.Join(pluginDir, fmt.Sprintf("plugin_%s_%s.module", runtime.GOOS, runtime.GOARCH))
	_, staterr := os.Stat(modPath)

	if staterr == nil {
		data, err := os.ReadFile(modPath)
		if err != nil {
			return nil, err
		}
		return wasmtime.NewModuleDeserialize(engine, data)
	}

	wmod, err := r.loadWASM(ctx, cache, expected)
	if err != nil {
		return nil, err
	}

	moduRegion := trace.StartRegion(ctx, "wasmtime.NewModule")
	module, err := wasmtime.NewModule(engine, wmod)
	moduRegion.End()
	if err != nil {
		return nil, fmt.Errorf("define wasi: %w", err)
	}

	if staterr != nil {
		// TODO: What permissions to use?
		err := os.MkdirAll(pluginDir, 0750)
		if err != nil && !os.IsExist(err) {
			return nil, fmt.Errorf("mkdirall: %w", err)
		}
		out, err := module.Serialize()
		if err != nil {
			return nil, fmt.Errorf("serialize: %w", err)
		}
		// TODO: What permissions to use?
		if err := os.WriteFile(modPath, out, 0666); err != nil {
			return nil, fmt.Errorf("cache wasm: %w", err)
		}
	}

	return module, nil
}

func (r *Runner) loadWASM(ctx context.Context, cache string, expected string) ([]byte, error) {
	pluginDir := filepath.Join(cache, expected)
	pluginPath := filepath.Join(pluginDir, "plugin.wasm")
	_, staterr := os.Stat(pluginPath)

	var body io.ReadCloser
	switch {
	case staterr == nil:
		file, err := os.Open(pluginPath)
		if err != nil {
			return nil, fmt.Errorf("os.Open: %s %w", pluginPath, err)
		}
		body = file

	case strings.HasPrefix(r.URL, "file://"):
		file, err := os.Open(strings.TrimPrefix(r.URL, "file://"))
		if err != nil {
			return nil, fmt.Errorf("os.Open: %s %w", r.URL, err)
		}
		body = file

	case strings.HasPrefix(r.URL, "https://"):
		// TODO: Set User-agent
		// TODO: Set ETag
		resp, err := http.Get(r.URL)
		if err != nil {
			return nil, fmt.Errorf("http.Get: %s %w", r.URL, err)
		}
		body = resp.Body

	default:
		return nil, fmt.Errorf("unknown scheme: %s", r.URL)
	}

	defer body.Close()

	wmod, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("readall: %w", err)
	}

	sum := sha256.Sum256(wmod)
	actual := fmt.Sprintf("%x", sum)

	if expected != actual {
		return nil, fmt.Errorf("invalid checksum: expected %s, got %s", expected, actual)
	}

	if staterr != nil {
		// TODO: What permissions to use?
		err := os.MkdirAll(pluginDir, 0750)
		if err != nil && !os.IsExist(err) {
			return nil, fmt.Errorf("mkdirall: %w", err)
		}
		// TODO: What permissions to use?
		if err := os.WriteFile(pluginPath, wmod, 0666); err != nil {
			return nil, fmt.Errorf("cache wasm: %w", err)
		}
	}

	return wmod, nil
}

func (r *Runner) Generate(ctx context.Context, req *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error) {
	stdinBlob, err := req.MarshalVT()
	if err != nil {
		return nil, err
	}

	engine := wasmtime.NewEngine()
	module, err := r.loadModule(ctx, engine)
	if err != nil {
		return nil, fmt.Errorf("loadModule: %w", err)
	}

	linker := wasmtime.NewLinker(engine)
	if err := linker.DefineWasi(); err != nil {
		return nil, err
	}

	dir, err := ioutil.TempDir("", "out")
	if err != nil {
		return nil, fmt.Errorf("temp dir: %w", err)
	}

	defer os.RemoveAll(dir)
	stdinPath := filepath.Join(dir, "stdin")
	stderrPath := filepath.Join(dir, "stderr")
	stdoutPath := filepath.Join(dir, "stdout")

	if err := os.WriteFile(stdinPath, stdinBlob, 0755); err != nil {
		return nil, fmt.Errorf("write file: %w", err)
	}

	// Configure WASI imports to write stdout into a file.
	wasiConfig := wasmtime.NewWasiConfig()
	wasiConfig.SetStdinFile(stdinPath)
	wasiConfig.SetStdoutFile(stdoutPath)
	wasiConfig.SetStderrFile(stderrPath)

	store := wasmtime.NewStore(engine)
	store.SetWasi(wasiConfig)

	linkRegion := trace.StartRegion(ctx, "linker.Instantiate")
	instance, err := linker.Instantiate(store, module)
	linkRegion.End()
	if err != nil {
		return nil, fmt.Errorf("define wasi: %w", err)
	}

	// Run the function
	callRegion := trace.StartRegion(ctx, "call _start")
	nom := instance.GetExport(store, "_start").Func()
	_, err = nom.Call(store)
	callRegion.End()
	if err != nil {
		return nil, fmt.Errorf("call: %w", err)
	}

	// Print WASM stdout
	stdoutBlob, err := os.ReadFile(stdoutPath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	var resp plugin.CodeGenResponse
	return &resp, resp.UnmarshalVT(stdoutBlob)
}
