package clickhouse

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestCatalogIntegration(t *testing.T) {
	schema := `
CREATE TABLE IF NOT EXISTS users
(
    id UInt32,
    name String,
    email String
)
ENGINE = MergeTree()
ORDER BY id;
`

	p := NewParser()
	stmts, err := p.Parse(strings.NewReader(schema))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(stmts))
	}

	// Debug: check what's in the statement
	if stmts[0].Raw != nil && stmts[0].Raw.Stmt != nil {
		if createStmt, ok := stmts[0].Raw.Stmt.(*ast.CreateTableStmt); ok {
			t.Logf("CreateTableStmt: Schema='%s', Table='%s'", createStmt.Name.Schema, createStmt.Name.Name)
			t.Logf("CreateTableStmt: Cols count=%d", len(createStmt.Cols))
		} else {
			t.Logf("Statement type: %T", stmts[0].Raw.Stmt)
		}
	}

	cat := NewCatalog()
	if cat.DefaultSchema != "default" {
		t.Errorf("Expected default schema 'default', got '%s'", cat.DefaultSchema)
	}

	// Try to update catalog with the CREATE TABLE
	t.Logf("Calling catalog.Update()...")
	err = cat.Update(stmts[0], nil)
	if err != nil {
		t.Fatalf("Catalog update failed: %v", err)
	}
	t.Logf("Catalog update succeeded")

	// Check if table was added
	t.Logf("Catalog has %d schemas", len(cat.Schemas))
	for i, schema := range cat.Schemas {
		t.Logf("Schema[%d]: Name='%s', Tables=%d", i, schema.Name, len(schema.Tables))
	}

	if len(cat.Schemas) == 0 {
		t.Fatal("No schemas in catalog")
	}

	defaultSchema := cat.Schemas[0]
	if len(defaultSchema.Tables) == 0 {
		t.Fatal("No tables in default schema")
	}

	table := defaultSchema.Tables[0]
	if table.Rel.Name != "users" {
		t.Errorf("Expected table name 'users', got '%s'", table.Rel.Name)
	}

	if len(table.Columns) != 3 {
		t.Errorf("Expected 3 columns, got %d", len(table.Columns))
	}

	// Log column types for debugging
	for i, col := range table.Columns {
		t.Logf("Column[%d]: Name='%s', Type.Name='%s', NotNull=%v", i, col.Name, col.Type.Name, col.IsNotNull)
	}
}
