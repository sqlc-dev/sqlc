package catalog_test

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// stubColumnGenerator satisfies the catalog's column generator dependency
// without pulling in the full compiler. CREATE VIEW only needs a column set
// to store; these tests assert on relation presence and dependency tracking,
// not on the view's column types.
type stubColumnGenerator struct{}

func (stubColumnGenerator) OutputColumns(ast.Node) ([]*catalog.Column, error) {
	return []*catalog.Column{
		{Name: "id", Type: ast.TypeName{Name: "int4"}},
	}, nil
}

func update(t *testing.T, c *catalog.Catalog, sql string) {
	t.Helper()
	stmts, err := postgresql.NewParser().Parse(strings.NewReader(sql))
	if err != nil {
		t.Fatal(err)
	}
	for _, stmt := range stmts {
		if err := c.Update(stmt, stubColumnGenerator{}); err != nil {
			t.Fatal(err)
		}
	}
}

func publicSchema(t *testing.T, c *catalog.Catalog) *catalog.Schema {
	t.Helper()
	for _, s := range c.Schemas {
		if s.Name == "public" {
			return s
		}
	}
	t.Fatal(`schema "public" not found`)
	return nil
}

func tableNames(schema *catalog.Schema) []string {
	names := make([]string, 0, len(schema.Tables))
	for _, tbl := range schema.Tables {
		names = append(names, tbl.Rel.Name)
	}
	return names
}

func TestDropTableCascadeEvictsDependentViews(t *testing.T) {
	c := catalog.New("public")
	update(t, c, `
CREATE TABLE base (id int);
CREATE VIEW child AS SELECT id FROM base;
CREATE VIEW grandchild AS SELECT id FROM child;
`)

	schema := publicSchema(t, c)
	if got := tableNames(schema); len(got) != 3 {
		t.Fatalf("expected 3 relations before drop, got %v", got)
	}

	update(t, c, `DROP TABLE base CASCADE;`)

	// base is dropped; child depends on base and grandchild depends on child,
	// so CASCADE must transitively evict both views.
	if got := tableNames(schema); len(got) != 0 {
		t.Fatalf("expected cascade drop to remove base and dependent views, got %v", got)
	}
}

func TestDropTableWithoutCascadeKeepsDependentViews(t *testing.T) {
	for _, tc := range []struct {
		name string
		sql  string
	}{
		{name: "restrict", sql: `DROP TABLE base RESTRICT;`},
		{name: "default", sql: `DROP TABLE base;`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := catalog.New("public")
			update(t, c, `
CREATE TABLE base (id int);
CREATE VIEW child AS SELECT id FROM base;
`)

			schema := publicSchema(t, c)
			update(t, c, tc.sql)

			got := tableNames(schema)
			if len(got) != 1 || got[0] != "child" {
				t.Fatalf("expected dependent view to remain after %s, got %v", tc.name, got)
			}
			if len(schema.Tables[0].DependsOn) != 1 || schema.Tables[0].DependsOn[0].Name != "base" {
				t.Fatalf("expected remaining view dependency to be preserved, got %#v", schema.Tables[0].DependsOn)
			}
		})
	}
}
