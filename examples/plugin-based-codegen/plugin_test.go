package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/sqlc-dev/sqlc/pkg/engine"
	"google.golang.org/protobuf/proto"
)

// TestEnginePlugin verifies that the SQLite3 engine plugin communicates correctly.
func TestEnginePlugin(t *testing.T) {
	ctx := context.Background()

	// Build the engine plugin
	pluginDir := filepath.Join("plugins", "sqlc-engine-sqlite3")
	pluginBin := filepath.Join(pluginDir, "sqlc-engine-sqlite3")

	buildCmd := exec.Command("go", "build", "-o", "sqlc-engine-sqlite3", ".")
	buildCmd.Dir = pluginDir
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build engine plugin: %v\n%s", err, output)
	}
	defer os.Remove(pluginBin)

	// Test Parse without schema
	t.Run("Parse", func(t *testing.T) {
		req := &engine.ParseRequest{
			Sql: "SELECT * FROM users WHERE id = ?;",
		}
		resp := &engine.ParseResponse{}
		if err := invokePlugin(ctx, pluginBin, "parse", req, resp); err != nil {
			t.Fatal(err)
		}

		if resp.Sql == "" {
			t.Fatal("expected non-empty sql in response")
		}
		if len(resp.Parameters) != 1 {
			t.Fatalf("expected 1 parameter, got %d", len(resp.Parameters))
		}
		t.Logf("✓ Parse: sql=%q params=%d columns=%d", resp.Sql, len(resp.Parameters), len(resp.Columns))
	})

	// Test Parse with schema (wildcard expansion)
	t.Run("ParseWithSchema", func(t *testing.T) {
		schemaSQL := `CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL
		);`
		req := &engine.ParseRequest{
			Sql:          "SELECT * FROM users WHERE id = ?;",
			SchemaSource: &engine.ParseRequest_SchemaSql{SchemaSql: schemaSQL},
		}
		resp := &engine.ParseResponse{}
		if err := invokePlugin(ctx, pluginBin, "parse", req, resp); err != nil {
			t.Fatal(err)
		}

		if resp.Sql == "" {
			t.Fatal("expected non-empty sql in response")
		}
		// With schema, wildcard should be expanded to explicit columns
		if len(resp.Columns) != 3 {
			t.Fatalf("expected 3 columns (id, name, email), got %d", len(resp.Columns))
		}
		if len(resp.Parameters) != 1 {
			t.Fatalf("expected 1 parameter, got %d", len(resp.Parameters))
		}
		t.Logf("✓ ParseWithSchema: sql=%q params=%d columns=%v", resp.Sql, len(resp.Parameters), columnNames(resp.Columns))
	})
}

func columnNames(cols []*engine.Column) []string {
	names := make([]string, len(cols))
	for i, c := range cols {
		names[i] = c.Name
	}
	return names
}

func invokePlugin(ctx context.Context, bin, method string, req, resp proto.Message) error {
	reqData, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, bin, method)
	cmd.Stdin = bytes.NewReader(reqData)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return proto.Unmarshal(stdout.Bytes(), resp)
}
