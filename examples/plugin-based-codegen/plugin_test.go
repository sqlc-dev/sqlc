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

	// Test Parse
	t.Run("Parse", func(t *testing.T) {
		req := &engine.ParseRequest{
			Sql: "SELECT * FROM users WHERE id = ?;",
		}
		resp := &engine.ParseResponse{}
		if err := invokePlugin(ctx, pluginBin, "parse", req, resp); err != nil {
			t.Fatal(err)
		}

		if len(resp.Statements) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(resp.Statements))
		}
		t.Logf("✓ Parse: %s", resp.Statements[0].RawSql)
	})

	// Test GetCatalog
	t.Run("GetCatalog", func(t *testing.T) {
		req := &engine.GetCatalogRequest{}
		resp := &engine.GetCatalogResponse{}
		if err := invokePlugin(ctx, pluginBin, "get_catalog", req, resp); err != nil {
			t.Fatal(err)
		}

		if resp.Catalog == nil || resp.Catalog.Name != "sqlite3" {
			t.Fatalf("expected catalog 'sqlite3', got %v", resp.Catalog)
		}
		t.Logf("✓ GetCatalog: %s (schema: %s)", resp.Catalog.Name, resp.Catalog.DefaultSchema)
	})

	// Test IsReservedKeyword
	t.Run("IsReservedKeyword", func(t *testing.T) {
		tests := []struct {
			keyword  string
			expected bool
		}{
			{"SELECT", true},
			{"PRAGMA", true},
			{"users", false},
		}

		for _, tc := range tests {
			req := &engine.IsReservedKeywordRequest{Keyword: tc.keyword}
			resp := &engine.IsReservedKeywordResponse{}
			if err := invokePlugin(ctx, pluginBin, "is_reserved_keyword", req, resp); err != nil {
				t.Fatal(err)
			}
			if resp.IsReserved != tc.expected {
				t.Errorf("IsReservedKeyword(%q) = %v, want %v", tc.keyword, resp.IsReserved, tc.expected)
			}
		}
		t.Log("✓ IsReservedKeyword")
	})

	// Test GetDialect
	t.Run("GetDialect", func(t *testing.T) {
		req := &engine.GetDialectRequest{}
		resp := &engine.GetDialectResponse{}
		if err := invokePlugin(ctx, pluginBin, "get_dialect", req, resp); err != nil {
			t.Fatal(err)
		}

		if resp.ParamStyle != "question" {
			t.Errorf("expected param_style 'question', got '%s'", resp.ParamStyle)
		}
		t.Logf("✓ GetDialect: quote=%s param=%s", resp.QuoteChar, resp.ParamStyle)
	})

	// Test GetCommentSyntax
	t.Run("GetCommentSyntax", func(t *testing.T) {
		req := &engine.GetCommentSyntaxRequest{}
		resp := &engine.GetCommentSyntaxResponse{}
		if err := invokePlugin(ctx, pluginBin, "get_comment_syntax", req, resp); err != nil {
			t.Fatal(err)
		}

		if !resp.Dash || !resp.SlashStar {
			t.Errorf("expected dash and slash_star comments")
		}
		t.Logf("✓ GetCommentSyntax: dash=%v slash_star=%v", resp.Dash, resp.SlashStar)
	})
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
