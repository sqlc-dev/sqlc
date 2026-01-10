package plugin

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

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
	"golang.org/x/sync/singleflight"

	"github.com/sqlc-dev/sqlc/internal/cache"
	"github.com/sqlc-dev/sqlc/internal/engine"
	"github.com/sqlc-dev/sqlc/internal/info"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

var wasmFlight singleflight.Group

type wasmRuntimeAndCode struct {
	rt   wazero.Runtime
	code wazero.CompiledModule
}

// WASMRunner runs an engine plugin as a WASM module.
type WASMRunner struct {
	URL    string
	SHA256 string
	Env    []string

	// Cached responses
	commentSyntax *WASMGetCommentSyntaxResponse
	dialect       *WASMGetDialectResponse
}

// NewWASMRunner creates a new WASMRunner.
func NewWASMRunner(url, sha256 string, env []string) *WASMRunner {
	return &WASMRunner{
		URL:    url,
		SHA256: sha256,
		Env:    env,
	}
}

func (r *WASMRunner) getChecksum(ctx context.Context) (string, error) {
	if r.SHA256 != "" {
		return r.SHA256, nil
	}
	_, sum, err := r.fetch(ctx, r.URL)
	if err != nil {
		return "", err
	}
	slog.Warn("fetching WASM binary to calculate sha256", "sha256", sum)
	return sum, nil
}

func (r *WASMRunner) fetch(ctx context.Context, uri string) ([]byte, string, error) {
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

func (r *WASMRunner) loadAndCompile(ctx context.Context) (*wasmRuntimeAndCode, error) {
	expected, err := r.getChecksum(ctx)
	if err != nil {
		return nil, err
	}

	cacheDir, err := cache.PluginsDir()
	if err != nil {
		return nil, err
	}

	value, err, _ := wasmFlight.Do(expected, func() (interface{}, error) {
		return r.loadAndCompileWASM(ctx, cacheDir, expected)
	})
	if err != nil {
		return nil, err
	}

	data, ok := value.(*wasmRuntimeAndCode)
	if !ok {
		return nil, fmt.Errorf("returned value was not a compiled module")
	}
	return data, nil
}

func (r *WASMRunner) loadAndCompileWASM(ctx context.Context, cacheDir string, expected string) (*wasmRuntimeAndCode, error) {
	pluginDir := filepath.Join(cacheDir, expected)
	pluginPath := filepath.Join(pluginDir, "engine.wasm")
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

	wazeroCache, err := wazero.NewCompilationCacheWithDir(filepath.Join(cacheDir, "wazero"))
	if err != nil {
		return nil, fmt.Errorf("wazero.NewCompilationCacheWithDir: %w", err)
	}

	config := wazero.NewRuntimeConfig().WithCompilationCache(wazeroCache)
	rt := wazero.NewRuntimeWithConfig(ctx, config)

	if _, err := wasi_snapshot_preview1.Instantiate(ctx, rt); err != nil {
		return nil, fmt.Errorf("wasi_snapshot_preview1 instantiate: %w", err)
	}

	code, err := rt.CompileModule(ctx, wmod)
	if err != nil {
		return nil, fmt.Errorf("compile module: %w", err)
	}

	return &wasmRuntimeAndCode{rt: rt, code: code}, nil
}

func (r *WASMRunner) invoke(ctx context.Context, method string, req, resp any) error {
	stdin, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}

	runtimeAndCode, err := r.loadAndCompile(ctx)
	if err != nil {
		return fmt.Errorf("loadBytes: %w", err)
	}

	var stderr, stdout bytes.Buffer

	conf := wazero.NewModuleConfig().
		WithName("").
		WithArgs("engine.wasm", method).
		WithStdin(bytes.NewReader(stdin)).
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
	if cerr := checkWASMError(err, stderr); cerr != nil {
		return cerr
	}

	if err := json.Unmarshal(stdout.Bytes(), resp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

func checkWASMError(err error, stderr bytes.Buffer) error {
	if err == nil {
		return err
	}

	if exitErr, ok := err.(*sys.ExitError); ok {
		if exitErr.ExitCode() == 0 {
			return nil
		}
	}

	stderrBlob := stderr.String()
	if len(stderrBlob) > 0 {
		return errors.New(stderrBlob)
	}
	return fmt.Errorf("call: %w", err)
}

// Parse implements engine.Parser.
func (r *WASMRunner) Parse(reader io.Reader) ([]ast.Statement, error) {
	sql, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	req := &WASMParseRequest{SQL: string(sql)}
	resp := &WASMParseResponse{}

	if err := r.invoke(context.Background(), "parse", req, resp); err != nil {
		return nil, err
	}

	var stmts []ast.Statement
	for _, s := range resp.Statements {
		node, err := parseASTJSON(s.ASTJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to parse AST: %w", err)
		}

		stmts = append(stmts, ast.Statement{
			Raw: &ast.RawStmt{
				Stmt:         node,
				StmtLocation: s.StmtLocation,
				StmtLen:      s.StmtLen,
			},
		})
	}

	return stmts, nil
}

// CommentSyntax implements engine.Parser.
func (r *WASMRunner) CommentSyntax() source.CommentSyntax {
	if r.commentSyntax == nil {
		req := &WASMGetCommentSyntaxRequest{}
		resp := &WASMGetCommentSyntaxResponse{}
		if err := r.invoke(context.Background(), "get_comment_syntax", req, resp); err != nil {
			return source.CommentSyntax{
				Dash:      true,
				SlashStar: true,
			}
		}
		r.commentSyntax = resp
	}

	return source.CommentSyntax{
		Dash:      r.commentSyntax.Dash,
		SlashStar: r.commentSyntax.SlashStar,
		Hash:      r.commentSyntax.Hash,
	}
}

// IsReservedKeyword implements engine.Parser.
func (r *WASMRunner) IsReservedKeyword(s string) bool {
	req := &WASMIsReservedKeywordRequest{Keyword: s}
	resp := &WASMIsReservedKeywordResponse{}
	if err := r.invoke(context.Background(), "is_reserved_keyword", req, resp); err != nil {
		return false
	}
	return resp.IsReserved
}

// GetCatalog returns the initial catalog for this engine.
func (r *WASMRunner) GetCatalog() (*catalog.Catalog, error) {
	req := &WASMGetCatalogRequest{}
	resp := &WASMGetCatalogResponse{}
	if err := r.invoke(context.Background(), "get_catalog", req, resp); err != nil {
		return nil, err
	}

	return convertWASMCatalog(&resp.Catalog), nil
}

// QuoteIdent implements engine.Dialect.
func (r *WASMRunner) QuoteIdent(s string) string {
	r.ensureDialect()
	if r.IsReservedKeyword(s) && r.dialect.QuoteChar != "" {
		return r.dialect.QuoteChar + s + r.dialect.QuoteChar
	}
	return s
}

// TypeName implements engine.Dialect.
func (r *WASMRunner) TypeName(ns, name string) string {
	if ns != "" {
		return ns + "." + name
	}
	return name
}

// Param implements engine.Dialect.
func (r *WASMRunner) Param(n int) string {
	r.ensureDialect()
	switch r.dialect.ParamStyle {
	case "dollar":
		return fmt.Sprintf("$%d", n)
	case "question":
		return "?"
	case "at":
		return fmt.Sprintf("@p%d", n)
	default:
		return fmt.Sprintf("$%d", n)
	}
}

// NamedParam implements engine.Dialect.
func (r *WASMRunner) NamedParam(name string) string {
	r.ensureDialect()
	if r.dialect.ParamPrefix != "" {
		return r.dialect.ParamPrefix + name
	}
	return "@" + name
}

// Cast implements engine.Dialect.
func (r *WASMRunner) Cast(arg, typeName string) string {
	r.ensureDialect()
	switch r.dialect.CastSyntax {
	case "double_colon":
		return arg + "::" + typeName
	default:
		return "CAST(" + arg + " AS " + typeName + ")"
	}
}

func (r *WASMRunner) ensureDialect() {
	if r.dialect == nil {
		req := &WASMGetDialectRequest{}
		resp := &WASMGetDialectResponse{}
		if err := r.invoke(context.Background(), "get_dialect", req, resp); err != nil {
			r.dialect = &WASMGetDialectResponse{
				QuoteChar:   `"`,
				ParamStyle:  "dollar",
				ParamPrefix: "@",
				CastSyntax:  "cast_function",
			}
		} else {
			r.dialect = resp
		}
	}
}

// convertWASMCatalog converts a WASM JSON Catalog to catalog.Catalog.
func convertWASMCatalog(c *WASMCatalog) *catalog.Catalog {
	if c == nil {
		return catalog.New("")
	}

	cat := catalog.New(c.DefaultSchema)
	cat.Name = c.Name
	cat.Comment = c.Comment
	cat.SearchPath = c.SearchPath

	cat.Schemas = make([]*catalog.Schema, 0, len(c.Schemas))
	for _, s := range c.Schemas {
		schema := &catalog.Schema{
			Name:    s.Name,
			Comment: s.Comment,
		}

		for _, t := range s.Tables {
			table := &catalog.Table{
				Rel: &ast.TableName{
					Catalog: t.Catalog,
					Schema:  t.Schema,
					Name:    t.Name,
				},
				Comment: t.Comment,
			}
			for _, col := range t.Columns {
				table.Columns = append(table.Columns, &catalog.Column{
					Name:       col.Name,
					Type:       ast.TypeName{Name: col.DataType},
					IsNotNull:  col.NotNull,
					IsArray:    col.IsArray,
					ArrayDims:  col.ArrayDims,
					Comment:    col.Comment,
					Length:     toPointerWASM(col.Length),
					IsUnsigned: col.IsUnsigned,
				})
			}
			schema.Tables = append(schema.Tables, table)
		}

		for _, e := range s.Enums {
			enum := &catalog.Enum{
				Name:    e.Name,
				Comment: e.Comment,
			}
			enum.Vals = append(enum.Vals, e.Values...)
			schema.Types = append(schema.Types, enum)
		}

		for _, f := range s.Functions {
			fn := &catalog.Function{
				Name:       f.Name,
				Comment:    f.Comment,
				ReturnType: &ast.TypeName{Schema: f.ReturnType.Schema, Name: f.ReturnType.Name},
			}
			for _, arg := range f.Args {
				fn.Args = append(fn.Args, &catalog.Argument{
					Name:       arg.Name,
					Type:       &ast.TypeName{Schema: arg.Type.Schema, Name: arg.Type.Name},
					HasDefault: arg.HasDefault,
				})
			}
			schema.Funcs = append(schema.Funcs, fn)
		}

		for _, t := range s.Types {
			schema.Types = append(schema.Types, &catalog.CompositeType{
				Name:    t.Name,
				Comment: t.Comment,
			})
		}

		cat.Schemas = append(cat.Schemas, schema)
	}

	return cat
}

func toPointerWASM(n int) *int {
	if n == 0 {
		return nil
	}
	return &n
}

// WASMPluginEngine wraps a WASMRunner to implement engine.Engine.
type WASMPluginEngine struct {
	name   string
	runner *WASMRunner
}

// NewWASMPluginEngine creates a new engine from a WASM plugin.
func NewWASMPluginEngine(name, url, sha256 string, env []string) *WASMPluginEngine {
	return &WASMPluginEngine{
		name:   name,
		runner: NewWASMRunner(url, sha256, env),
	}
}

// Name implements engine.Engine.
func (e *WASMPluginEngine) Name() string {
	return e.name
}

// Parser implements engine.Engine.
func (e *WASMPluginEngine) Parser() engine.Parser {
	return e.runner
}

// Catalog implements engine.Engine.
func (e *WASMPluginEngine) Catalog() *catalog.Catalog {
	cat, err := e.runner.GetCatalog()
	if err != nil {
		return catalog.New("")
	}
	return cat
}

// Selector implements engine.Engine.
func (e *WASMPluginEngine) Selector() engine.Selector {
	return &engine.DefaultSelector{}
}

// Dialect implements engine.Engine.
func (e *WASMPluginEngine) Dialect() engine.Dialect {
	return e.runner
}
