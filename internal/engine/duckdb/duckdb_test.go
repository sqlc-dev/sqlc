package duckdb_test

import (
	"testing"

	"github.com/sqlc-dev/sqlc/internal/compiler"
	"github.com/sqlc-dev/sqlc/internal/config"
)

func TestDuckDBEngineCreation(t *testing.T) {
	conf := config.SQL{
		Engine: config.EngineDuckDB,
	}

	combo := config.CombinedSettings{}

	c, err := compiler.NewCompiler(conf, combo)
	if err != nil {
		t.Fatalf("Failed to create DuckDB compiler: %v", err)
	}

	if c == nil {
		t.Fatal("Compiler is nil")
	}

	// Verify catalog was initialized
	catalog := c.Catalog()
	if catalog == nil {
		t.Fatal("Catalog is nil")
	}

	// Verify it uses PostgreSQL catalog (has pg_catalog schema)
	if catalog.DefaultSchema != "public" {
		t.Errorf("Expected default schema 'public', got %q", catalog.DefaultSchema)
	}
}
