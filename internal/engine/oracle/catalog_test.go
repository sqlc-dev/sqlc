package oracle

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func TestNewCatalog(t *testing.T) {
	c := NewCatalog()
	if c.DefaultSchema != defaultSchemaName {
		t.Fatalf("DefaultSchema = %q, want %q", c.DefaultSchema, defaultSchemaName)
	}
	if len(c.Schemas) != 1 {
		t.Fatalf("expected 1 schema, got %d", len(c.Schemas))
	}
	if len(c.Schemas[0].Funcs) == 0 {
		t.Fatalf("expected built-in functions to be seeded")
	}

	// Spot-check a few well-known Oracle functions are present.
	want := map[string]bool{"NVL": false, "TO_CHAR": false, "SYSDATE": false, "COUNT": false}
	for _, f := range c.Schemas[0].Funcs {
		if _, ok := want[f.Name]; ok {
			want[f.Name] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("expected built-in function %q to be registered", name)
		}
	}
}

func TestCatalogBuildCreateTable(t *testing.T) {
	p := NewParser()
	stmts, err := p.Parse(strings.NewReader(`CREATE TABLE employees (id NUMBER NOT NULL, name VARCHAR2(100))`))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	c := NewCatalog()
	if err := c.Build(stmts); err != nil {
		t.Fatalf("catalog.Build returned error: %v", err)
	}

	tbl, err := c.GetTable(&ast.TableName{Schema: defaultSchemaName, Name: "employees"})
	if err != nil {
		t.Fatalf("GetTable(employees) error: %v", err)
	}
	if tbl.Rel == nil || tbl.Rel.Name != "employees" {
		t.Fatalf("expected table 'employees', got %+v", tbl.Rel)
	}
	if len(tbl.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(tbl.Columns))
	}
	if tbl.Columns[0].Name != "id" {
		t.Errorf("col0 name = %q, want id", tbl.Columns[0].Name)
	}
	if !tbl.Columns[0].IsNotNull {
		t.Errorf("col0 (id) should be NOT NULL")
	}
}

func TestIsBuiltinType(t *testing.T) {
	for _, name := range []string{"NUMBER", "number", "VarChar2", "clob", "DATE", "TIMESTAMP"} {
		if !IsBuiltinType(name) {
			t.Errorf("IsBuiltinType(%q) = false, want true", name)
		}
	}
	for _, name := range []string{"jsonb", "geometry", ""} {
		if IsBuiltinType(name) {
			t.Errorf("IsBuiltinType(%q) = true, want false", name)
		}
	}
}

// ensure newTestCatalog is available and functional for future tests.
var _ = func() *catalog.Catalog { return newTestCatalog() }
