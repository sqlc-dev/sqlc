package ydb_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sqlc-dev/sqlc/internal/engine/ydb"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func TestCreateTable(t *testing.T) {
	tests := []struct {
		stmt string
		s    *catalog.Schema
	}{
		{
			stmt: `CREATE TABLE users (
				id Uint64,
				age Int32,
				score Float,
				PRIMARY KEY (id)
			)`,
			s: &catalog.Schema{
				Name: "main",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "users"},
						Columns: []*catalog.Column{
							{
								Name:      "id",
								Type:      ast.TypeName{Name: "Uint64"},
								IsNotNull: true,
							},
							{
								Name: "age",
								Type: ast.TypeName{Name: "Int32"},
							},
							{
								Name: "score",
								Type: ast.TypeName{Name: "Float"},
							},
						},
					},
				},
			},
		},
		{
			stmt: `CREATE TABLE posts (
				id Uint64,
				title Utf8 NOT NULL,
				content String,
				metadata Json,
				PRIMARY KEY (id)
			)`,
			s: &catalog.Schema{
				Name: "main",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "posts"},
						Columns: []*catalog.Column{
							{
								Name:      "id",
								Type:      ast.TypeName{Name: "Uint64"},
								IsNotNull: true,
							},
							{
								Name:      "title",
								Type:      ast.TypeName{Name: "Utf8"},
								IsNotNull: true,
							},
							{
								Name: "content",
								Type: ast.TypeName{Name: "String"},
							},
							{
								Name: "metadata",
								Type: ast.TypeName{Name: "Json"},
							},
						},
					},
				},
			},
		},
		{
			stmt: `CREATE TABLE orders (
				id Uuid,
				amount Decimal(22,9),
				created_at Uint64,
				PRIMARY KEY (id)
			)`,
			s: &catalog.Schema{
				Name: "main",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "orders"},
						Columns: []*catalog.Column{
							{
								Name:      "id",
								Type:      ast.TypeName{Name: "Uuid"},
								IsNotNull: true,
							},
							{
								Name: "amount",
								Type: ast.TypeName{
									Name: "decimal",
									Names: &ast.List{
										Items: []ast.Node{
											&ast.Integer{Ival: 22},
											&ast.Integer{Ival: 9},
										},
									},
								},
							},
							{
								Name: "created_at",
								Type: ast.TypeName{Name: "Uint64"},
							},
						},
					},
				},
			},
		},
	}

	p := ydb.NewParser()
	for i, tc := range tests {
		test := tc
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			stmts, err := p.Parse(strings.NewReader(test.stmt))
			if err != nil {
				t.Log(test.stmt)
				t.Fatal(err)
			}

			c := ydb.NewTestCatalog()
			if err := c.Build(stmts); err != nil {
				t.Log(test.stmt)
				t.Fatal(err)
			}

			e := ydb.NewTestCatalog()
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

			if diff := cmp.Diff(e, c, cmpopts.EquateEmpty(), cmpopts.IgnoreUnexported(catalog.Column{})); diff != "" {
				t.Log(test.stmt)
				t.Errorf("catalog mismatch:\n%s", diff)
			}
		})
	}
}
