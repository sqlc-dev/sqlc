package compiler

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/engine/clickhouse"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// TestClickHouseColumnNameResolution tests that column names are properly resolved
// from the catalog when processing SELECT statements with ClickHouse tables.
// This is a regression test for a bug where identifiers were being converted to
// string literals instead of ColumnRef nodes, preventing proper column name lookups.
func TestClickHouseColumnNameResolution(t *testing.T) {
	parser := clickhouse.NewParser()
	cat := clickhouse.NewCatalog()

	// Parse and add schema - CREATE DATABASE first, then table
	schemaSQL := `CREATE DATABASE IF NOT EXISTS sqlc_example;
CREATE TABLE sqlc_example.users (
	id UInt32,
	name String,
	email String,
	created_at DateTime
)`

	stmts, err := parser.Parse(strings.NewReader(schemaSQL))
	if err != nil {
		t.Fatalf("Parse schema failed: %v", err)
	}

	for _, stmt := range stmts {
		if err := cat.Update(stmt, nil); err != nil {
			t.Fatalf("Update catalog failed: %v", err)
		}
	}

	// Verify catalog is populated
	t.Logf("Catalog schemas: %d", len(cat.Schemas))
	found := false
	for _, schema := range cat.Schemas {
		if schema.Name == "sqlc_example" {
			found = true
			t.Logf("Found sqlc_example schema with %d tables", len(schema.Tables))
			if len(schema.Tables) > 0 {
				tbl := schema.Tables[0]
				t.Logf("  Table: %s.%s with %d columns", schema.Name, tbl.Rel.Name, len(tbl.Columns))
				for _, col := range tbl.Columns {
					t.Logf("    Column: %s (type: %v)", col.Name, col.Type)
				}
			}
		}
	}
	if !found {
		t.Fatal("sqlc_example schema not found in catalog")
	}

	// Create compiler
	conf := config.SQL{
		Engine: config.EngineClickHouse,
	}
	combo := config.CombinedSettings{
		Global: config.Config{},
	}

	c, err := NewCompiler(conf, combo)
	if err != nil {
		t.Fatalf("Failed to create compiler: %v", err)
	}

	// Replace the catalog with our populated one
	c.catalog = cat

	// Parse a SELECT query
	querySQL := "SELECT id, name, email FROM sqlc_example.users WHERE id = 1;"
	queryStmts, err := parser.Parse(strings.NewReader(querySQL))
	if err != nil {
		t.Fatalf("Parse query failed: %v", err)
	}

	if len(queryStmts) == 0 {
		t.Fatal("No queries parsed")
	}

	selectStmt := queryStmts[0].Raw.Stmt
	if selectStmt == nil {
		t.Fatal("Select statement is nil")
	}

	selectAst, ok := selectStmt.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Expected SelectStmt, got %T", selectStmt)
	}

	// Build the query catalog first
	qc, err := c.buildQueryCatalog(c.catalog, selectAst, nil)
	if err != nil {
		t.Fatalf("Failed to build query catalog: %v", err)
	}

	// Get output columns
	cols, err := c.outputColumns(qc, selectAst)
	if err != nil {
		t.Fatalf("Failed to get output columns: %v", err)
	}

	// Check if names are properly resolved
	if len(cols) != 3 {
		t.Errorf("Expected 3 columns, got %d", len(cols))
	}

	expectedNames := []string{"id", "name", "email"}
	for i, expected := range expectedNames {
		if i < len(cols) {
			if cols[i].Name != expected {
				t.Errorf("Column %d: expected name %q, got %q", i, expected, cols[i].Name)
			}
		}
	}
}
