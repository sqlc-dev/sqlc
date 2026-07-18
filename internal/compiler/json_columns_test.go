package compiler

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/sql/validate"
)

const testSchema = `
CREATE TABLE items (
	id bigserial PRIMARY KEY,
	x int NOT NULL,
	y text NOT NULL
);
`

func newTestCompiler(t *testing.T, schema string) *Compiler {
	t.Helper()
	parser := postgresql.NewParser()
	c := &Compiler{
		conf:    config.SQL{Engine: config.EnginePostgreSQL},
		parser:  parser,
		catalog: postgresql.NewCatalog(),
	}
	stmts, err := parser.Parse(strings.NewReader(schema))
	if err != nil {
		t.Fatalf("parse schema: %v", err)
	}
	for _, stmt := range stmts {
		if err := c.catalog.Update(stmt, c); err != nil {
			t.Fatalf("update catalog: %v", err)
		}
	}
	return c
}

func mustAnalyze(t *testing.T, c *Compiler, query string) *analysis {
	t.Helper()
	parser := postgresql.NewParser()
	stmts, err := parser.Parse(strings.NewReader(query))
	if err != nil {
		t.Fatalf("parse query: %v", err)
	}
	raw := stmts[0].Raw
	a, err := c.analyzeQuery(raw, query)
	if err != nil {
		t.Fatalf("analyze query: %v", err)
	}
	return a
}

func analyzeErr(t *testing.T, c *Compiler, query string) error {
	t.Helper()
	parser := postgresql.NewParser()
	stmts, err := parser.Parse(strings.NewReader(query))
	if err != nil {
		t.Fatalf("parse query: %v", err)
	}
	raw := stmts[0].Raw
	_, err = c.analyzeQuery(raw, query)
	if err == nil {
		t.Fatalf("expected an error, got none")
	}
	return err
}

func TestOutputColumnsJSON(t *testing.T) {
	c := newTestCompiler(t, testSchema)

	t.Run("scalar", func(t *testing.T) {
		a := mustAnalyze(t, c, `SELECT sqlc.jsonb_build_object."Obj"('x', x, 'y', y) AS obj FROM items;`)
		if len(a.Columns) != 1 {
			t.Fatalf("expected 1 column, got %d", len(a.Columns))
		}
		col := a.Columns[0]
		if col.JSONName != "Obj" {
			t.Errorf("JSONName = %q, want %q", col.JSONName, "Obj")
		}
		if col.IsArray {
			t.Errorf("IsArray = true, want false")
		}
		if !col.NotNull {
			t.Errorf("NotNull = false, want true")
		}
		if len(col.JSONFields) != 2 {
			t.Fatalf("expected 2 JSONFields, got %d", len(col.JSONFields))
		}
		if col.JSONFields[0].Name != "x" || col.JSONFields[0].DataType != "pg_catalog.int4" {
			t.Errorf("field 0 = %+v", col.JSONFields[0])
		}
		if col.JSONFields[1].Name != "y" || col.JSONFields[1].DataType != "text" {
			t.Errorf("field 1 = %+v", col.JSONFields[1])
		}
	})

	t.Run("array", func(t *testing.T) {
		a := mustAnalyze(t, c, `SELECT ARRAY(SELECT sqlc.jsonb_build_object."Obj"('x', x) FROM items) AS objs;`)
		col := a.Columns[0]
		if !col.IsArray || col.ArrayDims != 1 {
			t.Errorf("IsArray/ArrayDims = %v/%d, want true/1", col.IsArray, col.ArrayDims)
		}
		if col.JSONName != "Obj" {
			t.Errorf("JSONName = %q, want %q", col.JSONName, "Obj")
		}
	})

	t.Run("plain array subquery still types (not sqlc.json specific)", func(t *testing.T) {
		a := mustAnalyze(t, c, `SELECT ARRAY(SELECT x FROM items) AS xs;`)
		col := a.Columns[0]
		if !col.IsArray || col.ArrayDims != 1 {
			t.Errorf("IsArray/ArrayDims = %v/%d, want true/1", col.IsArray, col.ArrayDims)
		}
		if col.DataType != "pg_catalog.int4" {
			t.Errorf("DataType = %q, want %q", col.DataType, "pg_catalog.int4")
		}
	})

	t.Run("missing name", func(t *testing.T) {
		// sqlc.jsonb_build_object(...) requires a name; this is caught by
		// validate.SqlcFunctions before the query ever reaches outputColumns.
		parser := postgresql.NewParser()
		stmts, err := parser.Parse(strings.NewReader(`SELECT sqlc.jsonb_build_object('x', x) AS obj FROM items;`))
		if err != nil {
			t.Fatalf("parse query: %v", err)
		}
		err = validate.SqlcFunctions(stmts[0].Raw)
		if err == nil || !strings.Contains(err.Error(), "requires a name") {
			t.Errorf("error = %v, want mention of missing name", err)
		}
	})

	t.Run("odd args", func(t *testing.T) {
		err := analyzeErr(t, c, `SELECT sqlc.jsonb_build_object."Obj"('x', x, 'y') AS obj FROM items;`)
		if !strings.Contains(err.Error(), "even number") {
			t.Errorf("error = %v, want mention of even number of arguments", err)
		}
	})

	t.Run("non-literal key", func(t *testing.T) {
		err := analyzeErr(t, c, `SELECT sqlc.jsonb_build_object."Obj"(x, y) AS obj FROM items;`)
		if !strings.Contains(err.Error(), "string literal key") {
			t.Errorf("error = %v, want mention of string literal key", err)
		}
	})

	t.Run("array subquery with more than one column", func(t *testing.T) {
		err := analyzeErr(t, c, `SELECT ARRAY(SELECT x, y FROM items) AS objs;`)
		if !strings.Contains(err.Error(), "only one column") {
			t.Errorf("error = %v, want mention of one column", err)
		}
	})
}

func TestOutputColumnsEmbedJSON(t *testing.T) {
	c := newTestCompiler(t, testSchema)

	t.Run("scalar mirrors the table's columns", func(t *testing.T) {
		a := mustAnalyze(t, c, `SELECT sqlc.embed.jsonb(items) AS obj FROM items;`)
		col := a.Columns[0]
		if col.JSONName != "obj" {
			t.Errorf("JSONName = %q, want %q", col.JSONName, "obj")
		}
		if col.IsArray {
			t.Errorf("IsArray = true, want false")
		}
		if !col.NotNull {
			t.Errorf("NotNull = false, want true")
		}
		if len(col.JSONFields) != 3 {
			t.Fatalf("expected 3 JSONFields (id, x, y), got %d", len(col.JSONFields))
		}
		names := []string{col.JSONFields[0].Name, col.JSONFields[1].Name, col.JSONFields[2].Name}
		if names[0] != "id" || names[1] != "x" || names[2] != "y" {
			t.Errorf("field names = %v, want [id x y]", names)
		}
	})

	t.Run("array singularizes the outer alias for the element name", func(t *testing.T) {
		a := mustAnalyze(t, c, `SELECT ARRAY(SELECT sqlc.embed.jsonb(items) FROM items) AS objs;`)
		col := a.Columns[0]
		if !col.IsArray || col.ArrayDims != 1 {
			t.Errorf("IsArray/ArrayDims = %v/%d, want true/1", col.IsArray, col.ArrayDims)
		}
		if col.JSONName != "obj" {
			t.Errorf("JSONName = %q, want %q (singular of the alias)", col.JSONName, "obj")
		}
		if len(col.JSONFields) != 3 {
			t.Fatalf("expected 3 JSONFields, got %d", len(col.JSONFields))
		}
	})

	t.Run("unknown table", func(t *testing.T) {
		err := analyzeErr(t, c, `SELECT sqlc.embed.jsonb(missing) AS obj FROM items;`)
		if !strings.Contains(err.Error(), "table not found") {
			t.Errorf("error = %v, want mention of table not found", err)
		}
	})
}
