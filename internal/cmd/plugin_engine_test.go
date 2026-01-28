package cmd

// Engine-plugin pipeline integration tests.
//
// Why here (cmd) and not in endtoend?
//   - endtoend: black-box replay of full sqlc on testdata dirs, comparing stdout/stderr/output to golden files.
//   - These tests: unit-style integration of the engine-plugin path inside cmd — in-memory config and FileContents,
//     mocks for engine and codegen, no temp dirs or real plugins. They assert the data flow (schema+query → engine
//     mock → codegen request) and that the plugin package is used when PluginParseFunc is nil.
//
// Proof that the technology works:
//   - TestPluginPipeline_FullPipeline: one block → one Parse call; that call receives schema; codegen gets the result.
//   - TestPluginPipeline_NBlocksNCalls: N blocks in query.sql → exactly N Parse calls; each call receives schema.
//   - TestPluginPipeline_DatabaseOnly_ReceivesNoSchema: with analyzer.database: only + database.uri, each Parse
//     call receives empty schema (the real runner would get connection_params in ParseRequest).
//   - TestPluginPipeline_WithoutOverride_UsesPluginPackage: with PluginParseFunc nil, generate fails with an error
//     that is NOT "unknown engine", so we did enter runPluginQuerySet and call the engine process runner.

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/ext"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/plugin"
	pb "github.com/sqlc-dev/sqlc/pkg/engine"
)

const testPluginPipelineConfig = `version: "2"
sql:
  - engine: "testeng"
    schema: ["schema.sql"]
    queries: ["query.sql"]
    codegen:
      - plugin: "mock"
        out: "."
plugins:
  - name: "mock"
    process:
      cmd: "echo"
engines:
  - name: "testeng"
    process:
      cmd: "echo"
`

// engineMockRecord holds what the engine-plugin mock received and returned.
// Used to validate that the pipeline passes schema + raw query in, and that
// the plugin's Sql/Parameters/Columns are what later reach codegen.
// CalledWith records each Parse call so we can assert N blocks → N calls and
// that every call received schema (or "" in databaseOnly mode).
type engineMockRecord struct {
	Calls          int
	SchemaSQL      string // last call (backward compat)
	QuerySQL       string // last call (backward compat)
	CalledWith     []struct{ SchemaSQL, QuerySQL string }
	ReturnedSQL    string
	ReturnedParams []*pb.Parameter
	ReturnedCols   []*pb.Column
}

// codegenMockRecord holds what the codegen-plugin mock received.
type codegenMockRecord struct {
	Request *plugin.GenerateRequest
}

// TestPluginPipeline_FullPipeline validates the plugin-engine data flow end to end:
//
//  1. Inputs: schema and query file contents (from FileContents) are passed into the
//     engine-plugin mock. We assert exactly what the mock received.
//  2. Engine plugin returns: Sql, Parameters, Columns (the "enriched" result).
//  3. Pipeline converts that into compiler.Result and then into plugin.GenerateRequest.
//  4. Codegen plugin receives that request. We assert it contains the same SQL, params,
//     and columns, and that query name/cmd come from the query file comments.
//
// Note: the engine process runner is not called here because we use PluginParseFunc.
// That mock replaces the real ProcessRunner entirely. This test validates the cmd
// pipeline and the data contract at the boundaries; coverage of the plugin package
// comes from other tests (e.g. process runner, or an E2E test with a real engine binary).
func TestPluginPipeline_FullPipeline(t *testing.T) {
	ctx := context.Background()

	// --- Inputs (what we feed into the pipeline via FileContents) ---
	const (
		schemaContent = "CREATE TABLE users (id INT, name TEXT);"
		queryContent  = "-- name: GetUser :one\nSELECT id, name FROM users WHERE id = $1"
	)

	// --- Engine mock: record inputs, return "enriched" output ---
	engineRecord := &engineMockRecord{
		ReturnedSQL: "SELECT id, name FROM users WHERE id = $1",
		ReturnedParams: []*pb.Parameter{
			{Position: 1, DataType: "int", Nullable: false},
		},
		ReturnedCols: []*pb.Column{
			{Name: "id", DataType: "int", Nullable: false},
			{Name: "name", DataType: "text", Nullable: false},
		},
	}
	pluginParse := func(schemaSQL, querySQL string) (*pb.ParseResponse, error) {
		engineRecord.Calls++
		engineRecord.SchemaSQL = schemaSQL
		engineRecord.QuerySQL = querySQL
		engineRecord.CalledWith = append(engineRecord.CalledWith, struct{ SchemaSQL, QuerySQL string }{schemaSQL, querySQL})
		return &pb.ParseResponse{
			Sql:        engineRecord.ReturnedSQL,
			Parameters: engineRecord.ReturnedParams,
			Columns:    engineRecord.ReturnedCols,
		}, nil
	}

	// --- Codegen mock: record the full request ---
	codegenRecord := &codegenMockRecord{}
	mockCodegen := ext.HandleFunc(func(_ context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
		codegenRecord.Request = req
		return &plugin.GenerateResponse{}, nil
	})

	conf, err := config.ParseConfig(strings.NewReader(testPluginPipelineConfig))
	if err != nil {
		t.Fatalf("parse config: %v", err)
	}

	inputs := &sourceFiles{
		Config:     &conf,
		ConfigPath: "sqlc.yaml",
		Dir:        ".",
		FileContents: map[string][]byte{
			"schema.sql": []byte(schemaContent),
			"query.sql":  []byte(queryContent),
		},
	}

	var stderr bytes.Buffer
	debug := opts.DebugFromString("")
	debug.ProcessPlugins = true
	o := &Options{
		Env:                    Env{Debug: debug},
		Stderr:                 &stderr,
		PluginParseFunc:        pluginParse,
		CodegenHandlerOverride: mockCodegen,
	}

	_, err = generate(ctx, inputs, o)
	if err != nil {
		t.Fatalf("generate failed: %v\nstderr: %s", err, stderr.String())
	}

	// ---- 1) Validate what was sent INTO the engine plugin ----
	// N blocks in query.sql must yield N Parse calls; each call must receive schema (or connection in databaseOnly).
	if engineRecord.Calls != 1 {
		t.Errorf("engine mock called %d times, want 1 (one block → one Parse call)", engineRecord.Calls)
	}
	if len(engineRecord.CalledWith) != 1 {
		t.Errorf("engine mock CalledWith has %d entries, want 1", len(engineRecord.CalledWith))
	}
	if len(engineRecord.CalledWith) > 0 && engineRecord.CalledWith[0].SchemaSQL != schemaContent {
		t.Errorf("every Parse call must receive schema: got %q", engineRecord.CalledWith[0].SchemaSQL)
	}
	if engineRecord.SchemaSQL != schemaContent {
		t.Errorf("engine received schema:\n  got:  %q\n  want: %q", engineRecord.SchemaSQL, schemaContent)
	}
	// With one block, query SQL is the whole block (same as queryContent)
	if engineRecord.QuerySQL != queryContent {
		t.Errorf("engine received query:\n  got:  %q\n  want: %q", engineRecord.QuerySQL, queryContent)
	}

	// ---- 2) Validate what the codegen plugin received (must match engine output + metadata) ----
	if codegenRecord.Request == nil {
		t.Fatal("codegen mock was never called; request not recorded")
	}
	if len(codegenRecord.Request.Queries) == 0 {
		t.Fatal("codegen request has no queries")
	}
	q := codegenRecord.Request.Queries[0]

	// Name and Cmd come from the query file comment "-- name: GetUser :one"
	if got := q.GetName(); got != "GetUser" {
		t.Errorf("codegen query name = %q, want %q", got, "GetUser")
	}
	if got := q.GetCmd(); got != ":one" {
		t.Errorf("codegen query cmd = %q, want %q", got, ":one")
	}

	// Text must be exactly what the engine plugin returned (pipeline does not change it)
	if q.GetText() == "" {
		t.Error("codegen query has empty Text; plugin Sql did not reach codegen")
	}
	if got := q.GetText(); got != engineRecord.ReturnedSQL {
		t.Errorf("codegen query Text = %q, want (engine output) %q", got, engineRecord.ReturnedSQL)
	}

	// Params and columns must match what the engine plugin returned (codegen receives unchanged)
	if len(q.GetParams()) != len(engineRecord.ReturnedParams) {
		t.Errorf("codegen query has %d params, want %d", len(q.GetParams()), len(engineRecord.ReturnedParams))
	} else {
		for i, want := range engineRecord.ReturnedParams {
			p := q.GetParams()[i]
			if p.GetNumber() != want.Position {
				t.Errorf("param[%d] number = %d, want %d", i, p.GetNumber(), want.Position)
			}
			// plugin.Parameter.Column.Type.Name holds the data type
			if col := p.GetColumn(); col != nil && col.GetType() != nil {
				if got := col.GetType().GetName(); got != want.DataType {
					t.Errorf("param[%d] DataType = %q, want %q", i, got, want.DataType)
				}
			}
		}
	}
	if len(q.GetColumns()) != len(engineRecord.ReturnedCols) {
		t.Errorf("codegen query has %d columns, want %d", len(q.GetColumns()), len(engineRecord.ReturnedCols))
	} else {
		for i, want := range engineRecord.ReturnedCols {
			c := q.GetColumns()[i]
			if c.GetName() != want.Name {
				t.Errorf("column[%d] name = %q, want %q", i, c.GetName(), want.Name)
			}
			if typ := c.GetType(); typ != nil && typ.GetName() != want.DataType {
				t.Errorf("column[%d] type = %q, want %q", i, typ.GetName(), want.DataType)
			}
		}
	}

	// Sanity: codegen received exactly one query and we validated it
	if len(codegenRecord.Request.Queries) != 1 {
		t.Errorf("codegen received %d queries, expected 1", len(codegenRecord.Request.Queries))
	}
}

// TestPluginPipeline_WithoutOverride_UsesPluginPackage proves that when PluginParseFunc
// is not set, the pipeline calls the engine process runner (newEngineProcessRunner + parseRequest).
// It runs generate with a plugin engine and nil PluginParseFunc; we expect failure
// (e.g. from running "echo" as the engine binary), but the error must NOT be
// "unknown engine" — so we know we went past config lookup and into the plugin path.
// If you add panic("azaza") at the start of newEngineProcessRunner or parseRequest,
// this test will panic, confirming that the plugin package is actually invoked.
func TestPluginPipeline_WithoutOverride_UsesPluginPackage(t *testing.T) {
	ctx := context.Background()
	conf, err := config.ParseConfig(strings.NewReader(testPluginPipelineConfig))
	if err != nil {
		t.Fatalf("parse config: %v", err)
	}
	inputs := &sourceFiles{
		Config:     &conf,
		ConfigPath: "sqlc.yaml",
		Dir:        ".",
		FileContents: map[string][]byte{
			"schema.sql": []byte("CREATE TABLE t (id INT);"),
			"query.sql":  []byte("-- name: Get :one\nSELECT 1"),
		},
	}
	var stderr bytes.Buffer
	debug := opts.DebugFromString("")
	debug.ProcessPlugins = true
	o := &Options{
		Env:             Env{Debug: debug},
		Stderr:          &stderr,
		PluginParseFunc: nil, // do not override — must use built-in engine process runner
		CodegenHandlerOverride: ext.HandleFunc(func(_ context.Context, _ *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
			return &plugin.GenerateResponse{}, nil
		}),
	}

	_, err = generate(ctx, inputs, o)

	// We expect some error (e.g. "echo" does not speak the engine protocol).
	// What we must NOT see is "unknown engine" — that would mean we never reached
	// the plugin path. So the plugin package was used (ParseRequest or NewProcessRunner ran).
	if err == nil {
		t.Fatal("expected generate to fail when using real plugin runner with cmd=echo; nil error means plugin path was not exercised as intended")
	}
	if strings.Contains(err.Error(), "unknown engine") {
		t.Errorf("error is %q — we never entered the plugin path. With PluginParseFunc=nil, runPluginQuerySet must call the engine process runner.", err.Error())
	}
}

// TestPluginPipeline_NBlocksNCalls verifies that N sqlc blocks in query.sql yield N plugin Parse calls,
// and each call receives the schema (or connection params in databaseOnly mode).
func TestPluginPipeline_NBlocksNCalls(t *testing.T) {
	const (
		schemaContent = "CREATE TABLE users (id INT, name TEXT);"
		block1        = "-- name: GetUser :one\nSELECT id, name FROM users WHERE id = $1"
		block2        = "-- name: ListUsers :many\nSELECT id, name FROM users ORDER BY id"
	)
	queryContent := block1 + "\n\n" + block2
	// QueryBlocks slices from " name: " line to the next " name: " (exclusive), so first block includes "\n\n".
	expectedBlock1 := block1 + "\n\n"
	expectedBlock2 := block2

	engineRecord := &engineMockRecord{
		ReturnedSQL: "SELECT id, name FROM users WHERE id = $1",
		ReturnedParams: []*pb.Parameter{{Position: 1, DataType: "int", Nullable: false}},
		ReturnedCols:   []*pb.Column{{Name: "id", DataType: "int", Nullable: false}, {Name: "name", DataType: "text", Nullable: false}},
	}
	pluginParse := func(schemaSQL, querySQL string) (*pb.ParseResponse, error) {
		engineRecord.Calls++
		engineRecord.CalledWith = append(engineRecord.CalledWith, struct{ SchemaSQL, QuerySQL string }{schemaSQL, querySQL})
		return &pb.ParseResponse{Sql: querySQL, Parameters: engineRecord.ReturnedParams, Columns: engineRecord.ReturnedCols}, nil
	}
	conf, _ := config.ParseConfig(strings.NewReader(testPluginPipelineConfig))
	inputs := &sourceFiles{
		Config: &conf, ConfigPath: "sqlc.yaml", Dir: ".",
		FileContents: map[string][]byte{"schema.sql": []byte(schemaContent), "query.sql": []byte(queryContent)},
	}
	debug := opts.DebugFromString("")
	debug.ProcessPlugins = true
	o := &Options{
		Env: Env{Debug: debug}, Stderr: &bytes.Buffer{}, PluginParseFunc: pluginParse,
		CodegenHandlerOverride: ext.HandleFunc(func(_ context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) { return &plugin.GenerateResponse{}, nil }),
	}
	_, err := generate(context.Background(), inputs, o)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	if n := len(engineRecord.CalledWith); n != 2 {
		t.Errorf("expected 2 Parse calls (2 blocks), got %d", n)
	}
	for i, call := range engineRecord.CalledWith {
		if call.SchemaSQL != schemaContent {
			t.Errorf("Parse call %d: every call must receive schema; got schemaSQL %q", i+1, call.SchemaSQL)
		}
	}
	if len(engineRecord.CalledWith) >= 1 && engineRecord.CalledWith[0].QuerySQL != expectedBlock1 {
		t.Errorf("Parse call 1: query must be first block; got %q", engineRecord.CalledWith[0].QuerySQL)
	}
	if len(engineRecord.CalledWith) >= 2 && engineRecord.CalledWith[1].QuerySQL != expectedBlock2 {
		t.Errorf("Parse call 2: query must be second block; got %q", engineRecord.CalledWith[1].QuerySQL)
	}
}

const testPluginPipelineConfigDatabaseOnly = `version: "2"
sql:
  - engine: "testeng"
    schema: ["schema.sql"]
    queries: ["query.sql"]
    analyzer:
      database: only
    database:
      uri: "postgres://localhost/test"
    codegen:
      - plugin: "mock"
        out: "."
plugins:
  - name: "mock"
    process:
      cmd: "echo"
engines:
  - name: "testeng"
    process:
      cmd: "echo"
`

// TestPluginPipeline_DatabaseOnly_ReceivesNoSchema verifies that in databaseOnly mode (analyzer.database: only +
// database.uri) the plugin receives empty schema and the core would pass connection_params to the real runner.
// The mock only sees (schemaSQL, querySQL); in databaseOnly we pass schemaSQL="".
func TestPluginPipeline_DatabaseOnly_ReceivesNoSchema(t *testing.T) {
	const queryContent = "-- name: GetOne :one\nSELECT 1"
	engineRecord := &engineMockRecord{
		ReturnedSQL: "SELECT 1", ReturnedParams: nil, ReturnedCols: []*pb.Column{{Name: "?column?", DataType: "int", Nullable: true}},
	}
	pluginParse := func(schemaSQL, querySQL string) (*pb.ParseResponse, error) {
		engineRecord.Calls++
		engineRecord.CalledWith = append(engineRecord.CalledWith, struct{ SchemaSQL, QuerySQL string }{schemaSQL, querySQL})
		return &pb.ParseResponse{Sql: querySQL, Parameters: nil, Columns: engineRecord.ReturnedCols}, nil
	}
	conf, err := config.ParseConfig(strings.NewReader(testPluginPipelineConfigDatabaseOnly))
	if err != nil {
		t.Fatalf("parse config: %v", err)
	}
	inputs := &sourceFiles{
		Config: &conf, ConfigPath: "sqlc.yaml", Dir: ".",
		FileContents: map[string][]byte{"schema.sql": []byte("CREATE TABLE t (id INT);"), "query.sql": []byte(queryContent)},
	}
	debug := opts.DebugFromString("")
	debug.ProcessPlugins = true
	o := &Options{
		Env: Env{Debug: debug}, Stderr: &bytes.Buffer{}, PluginParseFunc: pluginParse,
		CodegenHandlerOverride: ext.HandleFunc(func(_ context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) { return &plugin.GenerateResponse{}, nil }),
	}
	_, err = generate(context.Background(), inputs, o)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	if len(engineRecord.CalledWith) != 1 {
		t.Errorf("expected 1 Parse call, got %d", len(engineRecord.CalledWith))
	}
	if len(engineRecord.CalledWith) > 0 && engineRecord.CalledWith[0].SchemaSQL != "" {
		t.Errorf("databaseOnly mode: each Parse call must receive empty schema (connection_params are used by real runner); got %q", engineRecord.CalledWith[0].SchemaSQL)
	}
	if len(engineRecord.CalledWith) > 0 && engineRecord.CalledWith[0].QuerySQL != queryContent {
		t.Errorf("query SQL must still be passed; got %q", engineRecord.CalledWith[0].QuerySQL)
	}
}

// TestPluginPipeline_OptionsOverrideNil ensures default Options do not inject mocks.
func TestPluginPipeline_OptionsOverrideNil(t *testing.T) {
	o := &Options{}
	if o.CodegenHandlerOverride != nil || o.PluginParseFunc != nil {
		t.Error("default Options should have nil overrides")
	}
}
