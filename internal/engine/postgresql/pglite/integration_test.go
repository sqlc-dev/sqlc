package pglite

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestAnalyzerIntegration(t *testing.T) {
	// Find the mock WASM file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not get current file path")
	}
	wasmPath := filepath.Join(filepath.Dir(filename), "testdata", "mock_pglite.wasm")

	// Check if WASM file exists
	wasmData, err := os.ReadFile(wasmPath)
	if err != nil {
		t.Skipf("mock_pglite.wasm not found (run 'GOOS=wasip1 GOARCH=wasm go build -o mock_pglite.wasm mock_pglite.go' in testdata/): %v", err)
	}

	// Calculate SHA256
	sum := sha256.Sum256(wasmData)
	sha := fmt.Sprintf("%x", sum)

	cfg := config.PGLite{
		URL:    "file://" + wasmPath,
		SHA256: sha,
	}

	analyzer := New(cfg)
	ctx := context.Background()
	defer analyzer.Close(ctx)

	migrations := []string{
		"CREATE TABLE users (id INTEGER NOT NULL, name TEXT, email TEXT NOT NULL, created_at TIMESTAMP)",
	}

	// Create a minimal AST node for position tracking
	node := &ast.TODO{}

	t.Run("simple select star", func(t *testing.T) {
		result, err := analyzer.Analyze(ctx, node, "SELECT * FROM users", migrations, nil)
		if err != nil {
			t.Fatalf("Analyze failed: %v", err)
		}

		if len(result.Columns) != 4 {
			t.Errorf("expected 4 columns, got %d", len(result.Columns))
		}

		// Check column names and types
		expectedCols := []struct {
			name    string
			notNull bool
		}{
			{"id", true},
			{"name", false},
			{"email", true},
			{"created_at", false},
		}

		for i, exp := range expectedCols {
			if i >= len(result.Columns) {
				break
			}
			col := result.Columns[i]
			if col.Name != exp.name {
				t.Errorf("column %d: expected name %q, got %q", i, exp.name, col.Name)
			}
			if col.NotNull != exp.notNull {
				t.Errorf("column %d (%s): expected NotNull=%v, got %v", i, exp.name, exp.notNull, col.NotNull)
			}
		}
	})

	t.Run("select with parameters", func(t *testing.T) {
		result, err := analyzer.Analyze(ctx, node, "SELECT id, name FROM users WHERE id = $1", migrations, nil)
		if err != nil {
			t.Fatalf("Analyze failed: %v", err)
		}

		if len(result.Columns) != 2 {
			t.Errorf("expected 2 columns, got %d", len(result.Columns))
		}

		if len(result.Params) != 1 {
			t.Errorf("expected 1 parameter, got %d", len(result.Params))
		}

		if len(result.Params) > 0 {
			if result.Params[0].Number != 1 {
				t.Errorf("expected param number 1, got %d", result.Params[0].Number)
			}
		}
	})

	t.Run("select specific columns", func(t *testing.T) {
		result, err := analyzer.Analyze(ctx, node, "SELECT id, email FROM users", migrations, nil)
		if err != nil {
			t.Fatalf("Analyze failed: %v", err)
		}

		if len(result.Columns) != 2 {
			t.Errorf("expected 2 columns, got %d", len(result.Columns))
		}

		if len(result.Columns) >= 2 {
			if result.Columns[0].Name != "id" {
				t.Errorf("expected first column 'id', got %q", result.Columns[0].Name)
			}
			if result.Columns[1].Name != "email" {
				t.Errorf("expected second column 'email', got %q", result.Columns[1].Name)
			}
		}
	})
}

func TestAnalyzerWithMultipleTables(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not get current file path")
	}
	wasmPath := filepath.Join(filepath.Dir(filename), "testdata", "mock_pglite.wasm")

	wasmData, err := os.ReadFile(wasmPath)
	if err != nil {
		t.Skipf("mock_pglite.wasm not found: %v", err)
	}

	sum := sha256.Sum256(wasmData)
	sha := fmt.Sprintf("%x", sum)

	cfg := config.PGLite{
		URL:    "file://" + wasmPath,
		SHA256: sha,
	}

	analyzer := New(cfg)
	ctx := context.Background()
	defer analyzer.Close(ctx)

	migrations := []string{
		"CREATE TABLE authors (id INTEGER NOT NULL, name TEXT NOT NULL)",
		"CREATE TABLE posts (id INTEGER NOT NULL, author_id INTEGER NOT NULL, title TEXT, body TEXT)",
	}

	node := &ast.TODO{}

	t.Run("query authors table", func(t *testing.T) {
		result, err := analyzer.Analyze(ctx, node, "SELECT * FROM authors", migrations, nil)
		if err != nil {
			t.Fatalf("Analyze failed: %v", err)
		}

		if len(result.Columns) != 2 {
			t.Errorf("expected 2 columns, got %d", len(result.Columns))
		}
	})

	t.Run("query posts table", func(t *testing.T) {
		result, err := analyzer.Analyze(ctx, node, "SELECT * FROM posts", migrations, nil)
		if err != nil {
			t.Fatalf("Analyze failed: %v", err)
		}

		if len(result.Columns) != 4 {
			t.Errorf("expected 4 columns, got %d", len(result.Columns))
		}
	})
}

func TestAnalyzerClose(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not get current file path")
	}
	wasmPath := filepath.Join(filepath.Dir(filename), "testdata", "mock_pglite.wasm")

	if _, err := os.Stat(wasmPath); os.IsNotExist(err) {
		t.Skipf("mock_pglite.wasm not found: %v", err)
	}

	wasmData, _ := os.ReadFile(wasmPath)
	sum := sha256.Sum256(wasmData)
	sha := fmt.Sprintf("%x", sum)

	cfg := config.PGLite{
		URL:    "file://" + wasmPath,
		SHA256: sha,
	}

	analyzer := New(cfg)
	ctx := context.Background()

	// Close without using should not error
	err := analyzer.Close(ctx)
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// Double close should not error
	err = analyzer.Close(ctx)
	if err != nil {
		t.Errorf("Double close failed: %v", err)
	}
}
