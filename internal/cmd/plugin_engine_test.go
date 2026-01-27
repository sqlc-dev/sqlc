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
//   - TestPluginPipeline_FullPipeline: with PluginParseFunc set, the pipeline sends schema+query into the mock,
//     takes Sql/Params/Columns from it, builds compiler.Result → plugin.GenerateRequest, and the codegen mock
//     receives that. So "plugin engine → ParseRequest contract → codegen" is validated.
//   - TestPluginPipeline_WithoutOverride_UsesPluginPackage: with PluginParseFunc nil, generate fails with an error
//     that is NOT "unknown engine", so we did enter runPluginQuerySet and call the engine process runner.
//     The plugin package is required for that path. Vet does not support plugin engines.

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
type engineMockRecord struct {
	Calls          int
	SchemaSQL      string
	QuerySQL       string
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
	if engineRecord.Calls != 1 {
		t.Errorf("engine mock called %d times, want 1", engineRecord.Calls)
	}
	if engineRecord.SchemaSQL != schemaContent {
		t.Errorf("engine received schema:\n  got:  %q\n  want: %q", engineRecord.SchemaSQL, schemaContent)
	}
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

// TestPluginPipeline_OptionsOverrideNil ensures default Options do not inject mocks.
func TestPluginPipeline_OptionsOverrideNil(t *testing.T) {
	o := &Options{}
	if o.CodegenHandlerOverride != nil || o.PluginParseFunc != nil {
		t.Error("default Options should have nil overrides")
	}
}
