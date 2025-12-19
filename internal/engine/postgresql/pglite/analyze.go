// Package pglite provides a PostgreSQL analyzer that uses PGLite running in WebAssembly.
// This allows for database-backed type analysis without requiring a running PostgreSQL server.
//
// To use this analyzer, enable it with SQLCEXPERIMENT=pglite and configure it in sqlc.yaml:
//
//	sql:
//	  - engine: postgresql
//	    analyzer:
//	      pglite:
//	        url: "file://path/to/pglite.wasm"
//	        sha256: "..."
package pglite

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"

	core "github.com/sqlc-dev/sqlc/internal/analysis"
	"github.com/sqlc-dev/sqlc/internal/cache"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/info"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/named"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

// Request types for communication with PGLite WASM module
type Request struct {
	Type       string   `json:"type"`       // "init", "exec", "prepare", "close"
	Migrations []string `json:"migrations"` // For "init": schema migrations to apply
	Query      string   `json:"query"`      // For "exec" and "prepare": SQL query
}

// Response from PGLite WASM module
type Response struct {
	Success bool            `json:"success"`
	Error   *ErrorResponse  `json:"error,omitempty"`
	Prepare *PrepareResult  `json:"prepare,omitempty"`
	Exec    *ExecResult     `json:"exec,omitempty"`
	Query   *QueryResult    `json:"query,omitempty"`
}

type ErrorResponse struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Position int    `json:"position"`
}

type PrepareResult struct {
	Columns []ColumnInfo    `json:"columns"`
	Params  []ParameterInfo `json:"params"`
}

type ColumnInfo struct {
	Name         string  `json:"name"`
	DataType     string  `json:"data_type"`
	DataTypeOID  uint32  `json:"data_type_oid"`
	NotNull      bool    `json:"not_null"`
	IsArray      bool    `json:"is_array"`
	ArrayDims    int     `json:"array_dims"`
	TableOID     uint32  `json:"table_oid,omitempty"`
	TableName    string  `json:"table_name,omitempty"`
	TableSchema  string  `json:"table_schema,omitempty"`
}

type ParameterInfo struct {
	Number      int    `json:"number"`
	DataType    string `json:"data_type"`
	DataTypeOID uint32 `json:"data_type_oid"`
	IsArray     bool   `json:"is_array"`
	ArrayDims   int    `json:"array_dims"`
}

type ExecResult struct {
	RowsAffected int64 `json:"rows_affected"`
}

type QueryResult struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}

// Analyzer implements the analyzer.Analyzer interface using PGLite WASM.
type Analyzer struct {
	cfg    config.PGLite
	mu     sync.Mutex
	rt     wazero.Runtime
	mod    api.Module
	inited bool
	schema []string

	// Caches for type lookups
	formats sync.Map
	columns sync.Map
	tables  sync.Map
}

// New creates a new PGLite analyzer with the given configuration.
func New(cfg config.PGLite) *Analyzer {
	return &Analyzer{
		cfg: cfg,
	}
}

// Analyze implements the analyzer.Analyzer interface.
// It prepares the given query against PGLite to extract column and parameter type information.
func (a *Analyzer) Analyze(ctx context.Context, n ast.Node, query string, migrations []string, ps *named.ParamSet) (*core.Analysis, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Initialize if not already done or if migrations changed
	if !a.inited || !equalMigrations(a.schema, migrations) {
		if err := a.init(ctx, migrations); err != nil {
			return nil, fmt.Errorf("pglite init: %w", err)
		}
		a.schema = migrations
		a.inited = true
	}

	// Prepare the query to get type information
	result, err := a.prepare(ctx, query)
	if err != nil {
		// Convert PGLite error to sqlerr.Error if possible
		var pgliteErr *PGLiteError
		if errors.As(err, &pgliteErr) {
			return nil, &sqlerr.Error{
				Code:     pgliteErr.Code,
				Message:  pgliteErr.Message,
				Location: max(n.Pos()+pgliteErr.Position-1, 0),
			}
		}
		return nil, err
	}

	var analysis core.Analysis

	// Convert columns
	for _, col := range result.Columns {
		dt := rewriteType(col.DataType)
		column := &core.Column{
			Name:         col.Name,
			OriginalName: col.Name,
			DataType:     dt,
			NotNull:      col.NotNull,
			IsArray:      col.IsArray,
			ArrayDims:    int32(col.ArrayDims),
		}
		if col.TableName != "" {
			column.Table = &core.Identifier{
				Schema: col.TableSchema,
				Name:   col.TableName,
			}
		}
		analysis.Columns = append(analysis.Columns, column)
	}

	// Convert parameters
	for _, param := range result.Params {
		dt := rewriteType(param.DataType)
		name := ""
		if ps != nil {
			name, _ = ps.NameFor(param.Number)
		}
		analysis.Params = append(analysis.Params, &core.Parameter{
			Number: int32(param.Number),
			Column: &core.Column{
				Name:      name,
				DataType:  dt,
				IsArray:   param.IsArray,
				ArrayDims: int32(param.ArrayDims),
				NotNull:   false, // Parameters are nullable by default
			},
		})
	}

	return &analysis, nil
}

// Close implements the analyzer.Analyzer interface.
func (a *Analyzer) Close(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.mod != nil {
		a.mod.Close(ctx)
		a.mod = nil
	}
	if a.rt != nil {
		a.rt.Close(ctx)
		a.rt = nil
	}
	a.inited = false
	return nil
}

// init initializes or reinitializes PGLite with the given migrations.
func (a *Analyzer) init(ctx context.Context, migrations []string) error {
	// Close existing runtime if any
	if a.rt != nil {
		if a.mod != nil {
			a.mod.Close(ctx)
		}
		a.rt.Close(ctx)
	}

	// Clear caches
	a.formats = sync.Map{}
	a.columns = sync.Map{}
	a.tables = sync.Map{}

	// Load and compile WASM
	wmod, err := a.loadWASM(ctx)
	if err != nil {
		return fmt.Errorf("load wasm: %w", err)
	}

	// Create wazero runtime with compilation cache
	cacheDir, err := cache.PluginsDir()
	if err != nil {
		return fmt.Errorf("cache dir: %w", err)
	}

	wazeroCache, err := wazero.NewCompilationCacheWithDir(filepath.Join(cacheDir, "pglite-wazero"))
	if err != nil {
		return fmt.Errorf("wazero cache: %w", err)
	}

	config := wazero.NewRuntimeConfig().WithCompilationCache(wazeroCache)
	a.rt = wazero.NewRuntimeWithConfig(ctx, config)

	// Instantiate WASI
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, a.rt); err != nil {
		return fmt.Errorf("wasi instantiate: %w", err)
	}

	// Compile and instantiate module
	compiled, err := a.rt.CompileModule(ctx, wmod)
	if err != nil {
		return fmt.Errorf("compile module: %w", err)
	}

	// Create request for initialization
	initReq := Request{
		Type:       "init",
		Migrations: migrations,
	}
	reqBytes, err := json.Marshal(initReq)
	if err != nil {
		return fmt.Errorf("marshal init request: %w", err)
	}

	var stdout, stderr bytes.Buffer

	modConfig := wazero.NewModuleConfig().
		WithName("pglite").
		WithArgs("pglite.wasm").
		WithStdin(bytes.NewReader(reqBytes)).
		WithStdout(&stdout).
		WithStderr(&stderr).
		WithEnv("SQLC_VERSION", info.Version).
		WithFSConfig(wazero.NewFSConfig())

	a.mod, err = a.rt.InstantiateModule(ctx, compiled, modConfig)
	if err != nil {
		errMsg := stderr.String()
		if errMsg != "" {
			return fmt.Errorf("instantiate module: %s", errMsg)
		}
		return fmt.Errorf("instantiate module: %w", err)
	}

	// Parse initialization response
	var resp Response
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		slog.Debug("pglite init response", "stdout", stdout.String(), "stderr", stderr.String())
		return fmt.Errorf("parse init response: %w", err)
	}

	if !resp.Success {
		if resp.Error != nil {
			return &PGLiteError{
				Code:    resp.Error.Code,
				Message: resp.Error.Message,
			}
		}
		return errors.New("pglite initialization failed")
	}

	return nil
}

// prepare sends a PREPARE request to PGLite and returns the result.
func (a *Analyzer) prepare(ctx context.Context, query string) (*PrepareResult, error) {
	req := Request{
		Type:  "prepare",
		Query: query,
	}

	resp, err := a.call(ctx, req)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		if resp.Error != nil {
			return nil, &PGLiteError{
				Code:     resp.Error.Code,
				Message:  resp.Error.Message,
				Position: resp.Error.Position,
			}
		}
		return nil, errors.New("prepare failed")
	}

	if resp.Prepare == nil {
		return nil, errors.New("prepare result missing")
	}

	return resp.Prepare, nil
}

// call sends a request to PGLite and returns the response.
// For a persistent module, this would use a different mechanism (e.g., function calls).
// For now, this demonstrates the interface that would be used.
func (a *Analyzer) call(ctx context.Context, req Request) (*Response, error) {
	// For modules that support function exports, we would call them directly.
	// Since PGLite WASM typically runs as a WASI command, we need to handle
	// persistent state differently.
	//
	// This implementation assumes the module exposes callable functions or
	// maintains state between invocations. In practice, you may need to:
	// 1. Use a module that exports query functions directly
	// 2. Re-instantiate with accumulated state
	// 3. Use a socket/pipe-based communication

	// Check if module has exported functions we can call
	queryFn := a.mod.ExportedFunction("pglite_query")
	if queryFn != nil {
		return a.callExported(ctx, queryFn, req)
	}

	// Fallback: re-instantiate with state (less efficient)
	return a.callViaReinstantiate(ctx, req)
}

// callExported calls an exported function on the WASM module.
func (a *Analyzer) callExported(ctx context.Context, fn api.Function, req Request) (*Response, error) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// Allocate memory for request
	malloc := a.mod.ExportedFunction("malloc")
	free := a.mod.ExportedFunction("free")

	if malloc == nil || free == nil {
		return nil, errors.New("module does not export malloc/free")
	}

	// Allocate memory for input
	results, err := malloc.Call(ctx, uint64(len(reqBytes)))
	if err != nil {
		return nil, fmt.Errorf("malloc: %w", err)
	}
	inputPtr := uint32(results[0])
	defer free.Call(ctx, uint64(inputPtr))

	// Write request to memory
	if !a.mod.Memory().Write(inputPtr, reqBytes) {
		return nil, errors.New("failed to write request to memory")
	}

	// Call the query function
	results, err = fn.Call(ctx, uint64(inputPtr), uint64(len(reqBytes)))
	if err != nil {
		return nil, fmt.Errorf("call: %w", err)
	}

	// Read response from memory
	outputPtr := uint32(results[0])
	outputLen := uint32(results[1])

	respBytes, ok := a.mod.Memory().Read(outputPtr, outputLen)
	if !ok {
		return nil, errors.New("failed to read response from memory")
	}

	var resp Response
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &resp, nil
}

// callViaReinstantiate handles modules that don't export callable functions.
// This re-instantiates the module with accumulated migrations plus the new query.
func (a *Analyzer) callViaReinstantiate(ctx context.Context, req Request) (*Response, error) {
	// For command-style WASM modules, we need to re-run them with the full state
	// This is less efficient but works with standard WASI command-line tools

	// Include migrations in the request so state is reconstructed
	fullReq := Request{
		Type:       req.Type,
		Query:      req.Query,
		Migrations: a.schema,
	}

	reqBytes, err := json.Marshal(fullReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// Load WASM again for fresh instance
	wmod, err := a.loadWASM(ctx)
	if err != nil {
		return nil, fmt.Errorf("load wasm: %w", err)
	}

	compiled, err := a.rt.CompileModule(ctx, wmod)
	if err != nil {
		return nil, fmt.Errorf("compile: %w", err)
	}

	var stdout, stderr bytes.Buffer

	modConfig := wazero.NewModuleConfig().
		WithName("").
		WithArgs("pglite.wasm", req.Type).
		WithStdin(bytes.NewReader(reqBytes)).
		WithStdout(&stdout).
		WithStderr(&stderr).
		WithEnv("SQLC_VERSION", info.Version)

	result, err := a.rt.InstantiateModule(ctx, compiled, modConfig)
	if err != nil {
		errMsg := stderr.String()
		if errMsg != "" {
			return nil, fmt.Errorf("instantiate: %s", errMsg)
		}
		return nil, fmt.Errorf("instantiate: %w", err)
	}
	defer result.Close(ctx)

	var resp Response
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		slog.Debug("pglite response", "stdout", stdout.String(), "stderr", stderr.String())
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &resp, nil
}

// loadWASM loads the PGLite WASM binary from the configured URL.
func (a *Analyzer) loadWASM(ctx context.Context) ([]byte, error) {
	url := a.cfg.URL
	expectedSHA := a.cfg.SHA256

	// Check cache first
	if expectedSHA != "" {
		cacheDir, err := cache.PluginsDir()
		if err == nil {
			cachePath := filepath.Join(cacheDir, expectedSHA, "pglite.wasm")
			if data, err := os.ReadFile(cachePath); err == nil {
				return data, nil
			}
		}
	}

	// Fetch the WASM binary
	data, actualSHA, err := fetch(ctx, url)
	if err != nil {
		return nil, err
	}

	// Verify checksum if provided
	if expectedSHA != "" && actualSHA != expectedSHA {
		return nil, fmt.Errorf("checksum mismatch: expected %s, got %s", expectedSHA, actualSHA)
	}

	// Warn if no checksum provided
	if expectedSHA == "" {
		slog.Warn("pglite: no sha256 checksum provided, set sha256 in config for security", "actual_sha256", actualSHA)
	}

	// Cache the binary
	if expectedSHA != "" {
		cacheDir, err := cache.PluginsDir()
		if err == nil {
			pluginDir := filepath.Join(cacheDir, expectedSHA)
			if err := os.MkdirAll(pluginDir, 0755); err == nil {
				os.WriteFile(filepath.Join(pluginDir, "pglite.wasm"), data, 0444)
			}
		}
	}

	return data, nil
}

// fetch downloads content from a URL (file:// or https://).
func fetch(ctx context.Context, url string) ([]byte, string, error) {
	var body io.ReadCloser

	switch {
	case strings.HasPrefix(url, "file://"):
		path := strings.TrimPrefix(url, "file://")
		file, err := os.Open(path)
		if err != nil {
			return nil, "", fmt.Errorf("open file: %w", err)
		}
		body = file

	case strings.HasPrefix(url, "https://"):
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, "", fmt.Errorf("create request: %w", err)
		}
		req.Header.Set("User-Agent", fmt.Sprintf("sqlc/%s Go/%s (%s %s)", info.Version, runtime.Version(), runtime.GOOS, runtime.GOARCH))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, "", fmt.Errorf("fetch: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, "", fmt.Errorf("fetch failed: %s", resp.Status)
		}
		body = resp.Body

	default:
		return nil, "", fmt.Errorf("unsupported URL scheme: %s", url)
	}

	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, "", fmt.Errorf("read: %w", err)
	}

	sum := sha256.Sum256(data)
	checksum := fmt.Sprintf("%x", sum)

	return data, checksum, nil
}

// PGLiteError represents an error from PGLite.
type PGLiteError struct {
	Code     string
	Message  string
	Position int
}

func (e *PGLiteError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return e.Message
}

// rewriteType converts PostgreSQL type names to the canonical form used by sqlc.
func rewriteType(dt string) string {
	switch {
	case strings.HasPrefix(dt, "character("):
		return "pg_catalog.bpchar"
	case strings.HasPrefix(dt, "character varying"):
		return "pg_catalog.varchar"
	case strings.HasPrefix(dt, "bit varying"):
		return "pg_catalog.varbit"
	case strings.HasPrefix(dt, "bit("):
		return "pg_catalog.bit"
	}
	switch dt {
	case "bpchar":
		return "pg_catalog.bpchar"
	case "timestamp without time zone":
		return "pg_catalog.timestamp"
	case "timestamp with time zone":
		return "pg_catalog.timestamptz"
	case "time without time zone":
		return "pg_catalog.time"
	case "time with time zone":
		return "pg_catalog.timetz"
	}
	return dt
}

// equalMigrations compares two migration slices for equality.
func equalMigrations(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
