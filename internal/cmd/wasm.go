package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/trace"

	wasmtime "github.com/bytecodealliance/wasmtime-go"

	"github.com/kyleconroy/sqlc/internal/plugin"
)

//go:embed sqlc-codegen-python.wasm
var pythonCodeGen []byte

func pythonGenerate(req *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error) {
	stdinBlob, err := req.MarshalVT()
	if err != nil {
		return nil, err
	}

	cctx := context.Background()

	ctx, task := trace.NewTask(cctx, "pythonGenerate")
	defer task.End()

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

	engine := wasmtime.NewEngine()
	linker := wasmtime.NewLinker(engine)

	// Link WASI
	if err := linker.DefineWasi(); err != nil {
		return nil, fmt.Errorf("define wasi: %w", err)
	}

	// Configure WASI imports to write stdout into a file.
	wasiConfig := wasmtime.NewWasiConfig()
	wasiConfig.SetStdinFile(stdinPath)
	wasiConfig.SetStdoutFile(stdoutPath)
	wasiConfig.SetStderrFile(stderrPath)

	store := wasmtime.NewStore(engine)
	store.SetWasi(wasiConfig)

	// Set the version to the same as in the WAT.
	// wasi, err := wasmtime.NewWasiInstance(store, wasiConfig, "wasi_snapshot_preview1")
	// if err != nil {
	// 	return fmt.Errorf("new wasi instances: %w", err)
	// }

	// Create our module
	//
	// Compiling modules requires WebAssembly binary input, but the wasmtime
	// package also supports converting the WebAssembly text format to the
	// binary format.
	// wasm, err := os.ReadFile("plugin.wasm")
	// if err != nil {
	// 	return fmt.Errorf("read file: %w", err)
	// }

	moduRegion := trace.StartRegion(ctx, "wasmtime.NewModule")
	module, err := wasmtime.NewModule(store.Engine, pythonCodeGen)
	moduRegion.End()
	if err != nil {
		return nil, fmt.Errorf("define wasi: %w", err)
	}

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
