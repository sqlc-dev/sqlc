package postgresql

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"

	"github.com/google/go-cmp/cmp"
)

func TestUpdateErrors(t *testing.T) {
	p := NewParser()
	for i, tc := range []struct {
		stmt string
		err  *sqlerr.Error
	}{
		{
			`
			CREATE TABLE foo ();
			CREATE TABLE foo ();
			`,
			sqlerr.RelationExists("foo"),
		},
		{
			`
			CREATE TYPE foo AS ENUM ('bar');
			CREATE TYPE foo AS ENUM ('bar');
			`,
			sqlerr.TypeExists("foo"),
		},
		{
			`
			DROP TABLE foo;
			`,
			sqlerr.RelationNotFound("foo"),
		},
		{
			`
			DROP TYPE foo;
			`,
			sqlerr.TypeNotFound("foo"),
		},
		{
			`
			CREATE TABLE foo ();
			CREATE TABLE bar ();
			ALTER TABLE foo RENAME TO bar;
			`,
			sqlerr.RelationExists("bar"),
		},
		{
			`
			ALTER TABLE foo RENAME TO bar;
			`,
			sqlerr.RelationNotFound("foo"),
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo ADD COLUMN bar text;
			ALTER TABLE foo ADD COLUMN bar text;
			`,
			sqlerr.ColumnExists("foo", "bar"),
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo DROP COLUMN bar;
			`,
			sqlerr.ColumnNotFound("foo", "bar"),
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo ALTER COLUMN bar SET NOT NULL;
			`,
			sqlerr.ColumnNotFound("foo", "bar"),
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo ALTER COLUMN bar DROP NOT NULL;
			`,
			sqlerr.ColumnNotFound("foo", "bar"),
		},
		{
			`
			CREATE SCHEMA foo;
			CREATE SCHEMA foo;
			`,
			sqlerr.SchemaExists("foo"),
		},
		{
			`
			ALTER TABLE foo.baz SET SCHEMA bar;
			`,
			sqlerr.SchemaNotFound("foo"),
		},
		{
			`
			CREATE SCHEMA foo;
			ALTER TABLE foo.baz SET SCHEMA bar;
			`,
			sqlerr.RelationNotFound("baz"),
		},
		{
			`
			CREATE SCHEMA foo;
			CREATE TABLE foo.baz ();
			ALTER TABLE foo.baz SET SCHEMA bar;
			`,
			sqlerr.SchemaNotFound("bar"),
		},
		{
			`
			DROP SCHEMA bar;
			`,
			sqlerr.SchemaNotFound("bar"),
		},
		{
			`
			ALTER TABLE foo RENAME bar TO baz;
			`,
			sqlerr.RelationNotFound("foo"),
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo RENAME bar TO baz;
			`,
			sqlerr.ColumnNotFound("foo", "bar"),
		},
		{
			`
			CREATE TABLE foo (bar text, baz text);
			ALTER TABLE foo RENAME bar TO baz;
			`,
			sqlerr.ColumnExists("foo", "baz"),
		},
	} {
		test := tc
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			stmts, err := p.Parse(strings.NewReader(test.stmt))
			if err != nil {
				t.Log(test.stmt)
				t.Fatal(err)
			}

			c := NewCatalog()
			err = c.Build(stmts)
			if err == nil {
				t.Log(test.stmt)
				t.Fatal("err was nil")
			}

			var actual *sqlerr.Error
			if !errors.As(err, &actual) {
				t.Fatalf("err is not *sqlerr.Error: %#v", err)
			}

			if diff := cmp.Diff(test.err.Error(), actual.Error()); diff != "" {
				t.Log(test.stmt)
				t.Errorf("error mismatch: \n%s", diff)
			}
		})
	}
}

func TestDropTableCascadeViewRecreate(t *testing.T) {
	// Regression test for https://github.com/sqlc-dev/sqlc/issues/4416
	// DROP TABLE CASCADE should remove dependent views from the catalog,
	// allowing a subsequent CREATE VIEW with the same name to succeed.
	p := NewParser()

	// First: create the table
	stmts1, err := p.Parse(strings.NewReader(`
		CREATE TABLE reference_rates (id BIGSERIAL PRIMARY KEY);
	`))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	c := NewCatalog()
	if err := c.Build(stmts1); err != nil {
		t.Fatalf("create table error: %v", err)
	}

	// Manually add a view that depends on reference_rates to the catalog
	var schema *catalog.Schema
	for _, s := range c.Schemas {
		if s.Name == "public" {
			schema = s
		}
	}
	schema.Tables = append(schema.Tables, &catalog.Table{
		Rel:     &ast.TableName{Schema: "public", Name: "vw_reference_rates"},
		Columns: []*catalog.Column{{Name: "id"}},
		DependsOnTables: []*ast.TableName{
			{Schema: "public", Name: "reference_rates"},
		},
	})

	// Verify the view exists
	if !viewExists(schema, "vw_reference_rates") {
		t.Fatal("view not found in catalog before drop")
	}

	// Second: DROP TABLE CASCADE
	stmts2, err := p.Parse(strings.NewReader(`
		DROP TABLE reference_rates CASCADE;
	`))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := c.Build(stmts2); err != nil {
		t.Fatalf("DROP TABLE CASCADE error: %v", err)
	}

	// Verify the view was removed
	if viewExists(schema, "vw_reference_rates") {
		t.Fatal("expected view to be removed by CASCADE, but it still exists")
	}
}

func TestDropTableCascadeWithoutCascadeFails(t *testing.T) {
	// Without CASCADE, dropping a table that has a dependent view leaves the view
	// in the catalog (matching current sqlc behavior, though real PostgreSQL would
	// reject DROP TABLE without CASCADE when views depend on it).
	p := NewParser()

	// Create the table
	stmts1, err := p.Parse(strings.NewReader(`
		CREATE TABLE reference_rates (id BIGSERIAL PRIMARY KEY);
	`))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	c := NewCatalog()
	if err := c.Build(stmts1); err != nil {
		t.Fatalf("create table error: %v", err)
	}

	// Manually add a view that depends on reference_rates
	schema := c.Schemas[0]
	for _, s := range c.Schemas {
		if s.Name == "public" {
			schema = s
		}
	}
	schema.Tables = append(schema.Tables, &catalog.Table{
		Rel:     &ast.TableName{Schema: "public", Name: "vw_reference_rates"},
		Columns: []*catalog.Column{{Name: "id"}},
		DependsOnTables: []*ast.TableName{
			{Schema: "public", Name: "reference_rates"},
		},
	})

	// DROP TABLE without CASCADE
	stmts2, err := p.Parse(strings.NewReader(`
		DROP TABLE reference_rates;
	`))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := c.Build(stmts2); err != nil {
		t.Fatalf("DROP TABLE error: %v", err)
	}

	// Without CASCADE, the view should still exist in the catalog
	if !viewExists(schema, "vw_reference_rates") {
		t.Fatal("expected view to still exist without CASCADE, but it was removed")
	}
}

func viewExists(schema *catalog.Schema, name string) bool {
	for _, tbl := range schema.Tables {
		if tbl.Rel.Name == name {
			return true
		}
	}
	return false
}
