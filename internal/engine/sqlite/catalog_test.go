package sqlite

import (
	"strconv"
	"strings"
	"testing"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestUpdate(t *testing.T) {
	p := NewParser()

	for i, tc := range []struct {
		stmt string
		s    *catalog.Schema
	}{
		{
			`
			CREATE TABLE foo (bar text);
			`,
			&catalog.Schema{
				Name: "main",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "foo"},
						Columns: []*catalog.Column{
							{
								Name: "bar",
								Type: ast.TypeName{Name: "text"},
							},
						},
					},
				},
			},
		},
		{
			`
			CREATE TABLE foo (bar text);
			ALTER TABLE foo RENAME TO baz;
			`,
			&catalog.Schema{
				Name: "main",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "baz"},
						Columns: []*catalog.Column{
							{
								Name: "bar",
								Type: ast.TypeName{Name: "text"},
							},
						},
					},
				},
			},
		},
		{
			`
			CREATE TABLE foo (bar text);
			ALTER TABLE foo ADD COLUMN baz bool;
			`,
			&catalog.Schema{
				Name: "main",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "foo"},
						Columns: []*catalog.Column{
							{
								Name: "bar",
								Type: ast.TypeName{Name: "text"},
							},
							{
								Name: "baz",
								Type: ast.TypeName{Name: "bool"},
							},
						},
					},
				},
			},
		},
		{
			`
			CREATE TABLE foo (bar text);
			ALTER TABLE foo RENAME COLUMN bar TO baz;
			`,
			&catalog.Schema{
				Name: "main",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "foo"},
						Columns: []*catalog.Column{
							{
								Name: "baz",
								Type: ast.TypeName{Name: "text"},
							},
						},
					},
				},
			},
		},
		{
			`
			CREATE TABLE foo (bar text);
			ALTER TABLE foo RENAME bar TO baz;
			`,
			&catalog.Schema{
				Name: "main",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "foo"},
						Columns: []*catalog.Column{
							{
								Name: "baz",
								Type: ast.TypeName{Name: "text"},
							},
						},
					},
				},
			},
		},
		{
			`
			ATTACH ':memory:' as ns;
			CREATE TABLE ns.foo (bar text);
			`,
			&catalog.Schema{
				Name: "ns",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Schema: "ns", Name: "foo"},
						Columns: []*catalog.Column{
							{
								Name: "bar",
								Type: ast.TypeName{Name: "text"},
							},
						},
					},
				},
			},
		},
		{
			`
			ATTACH ':memory:' as ns;
			CREATE TABLE ns.foo (bar text);
			ALTER TABLE ns.foo RENAME TO baz;
			`,
			&catalog.Schema{
				Name: "ns",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Schema: "ns", Name: "baz"},
						Columns: []*catalog.Column{
							{
								Name: "bar",
								Type: ast.TypeName{Name: "text"},
							},
						},
					},
				},
			},
		},
		{
			`
			ATTACH ':memory:' as ns;
			CREATE TABLE ns.foo (bar text);
			ALTER TABLE ns.foo ADD COLUMN baz bool;
			`,
			&catalog.Schema{
				Name: "ns",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Schema: "ns", Name: "foo"},
						Columns: []*catalog.Column{
							{
								Name: "bar",
								Type: ast.TypeName{Name: "text"},
							},
							{
								Name: "baz",
								Type: ast.TypeName{Name: "bool"},
							},
						},
					},
				},
			},
		},
		{
			`
			ATTACH ':memory:' as ns;
			CREATE TABLE ns.foo (bar text);
			ALTER TABLE ns.foo RENAME COLUMN bar TO baz;
			`,
			&catalog.Schema{
				Name: "ns",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Schema: "ns", Name: "foo"},
						Columns: []*catalog.Column{
							{
								Name: "baz",
								Type: ast.TypeName{Name: "text"},
							},
						},
					},
				},
			},
		},
		{
			`
			ATTACH ':memory:' as ns;
			CREATE TABLE ns.foo (bar text);
			ALTER TABLE ns.foo RENAME bar TO baz;
			`,
			&catalog.Schema{
				Name: "ns",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Schema: "ns", Name: "foo"},
						Columns: []*catalog.Column{
							{
								Name: "baz",
								Type: ast.TypeName{Name: "text"},
							},
						},
					},
				},
			},
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
			if err := c.Build(stmts); err != nil {
				t.Log(test.stmt)
				t.Fatal(err)
			}

			e := NewCatalog()
			if test.s != nil {
				var replaced bool
				for i := range e.Schemas {
					if e.Schemas[i].Name == test.s.Name {
						e.Schemas[i] = test.s
						replaced = true
						break
					}
				}
				if !replaced {
					e.Schemas = append(e.Schemas, test.s)
				}
			}

			if diff := cmp.Diff(e, c, cmpopts.EquateEmpty()); diff != "" {
				t.Log(test.stmt)
				t.Errorf("catalog mismatch:\n%s", diff)
			}
		})
	}
}
