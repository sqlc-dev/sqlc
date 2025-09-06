package ydb_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sqlc-dev/sqlc/internal/engine/ydb"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func strPtr(s string) *string {
	return &s
}

func TestSelect(t *testing.T) {
	tests := []struct {
		stmt     string
		expected ast.Node
	}{
		// Basic Types Select
		{
			stmt: `SELECT 52`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.A_Const{
										Val: &ast.Integer{Ival: 52},
									},
								},
							},
						},
						FromClause: &ast.List{},
					},
				},
			},
		},
		{
			stmt: `SELECT 'hello'`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.A_Const{
										Val: &ast.String{Str: "hello"},
									},
								},
							},
						},
						FromClause: &ast.List{},
					},
				},
			},
		},
		{
			stmt: `SELECT 'it\'s string with quote in it'`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.A_Const{
										Val: &ast.String{Str: `it\'s string with quote in it`},
									},
								},
							},
						},
						FromClause: &ast.List{},
					},
				},
			},
		},
		{
			stmt: "SELECT 3.14",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.A_Const{
										Val: &ast.Float{Str: "3.14"},
									},
								},
							},
						},
						FromClause: &ast.List{},
					},
				},
			},
		},
		{
			stmt: "SELECT NULL",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.Null{},
								},
							},
						},
						FromClause: &ast.List{},
					},
				},
			},
		},
		{
			stmt: "SELECT true",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.A_Const{
										Val: &ast.Boolean{Boolval: true},
									},
								},
							},
						},
						FromClause: &ast.List{},
					},
				},
			},
		},
		{
			stmt: "SELECT 2+3*4",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.A_Expr{
										Name: &ast.List{
											Items: []ast.Node{
												&ast.String{Str: "+"},
											},
										},
										Lexpr: &ast.A_Const{
											Val: &ast.Integer{Ival: 2},
										},
										Rexpr: &ast.A_Expr{
											Name: &ast.List{
												Items: []ast.Node{
													&ast.String{Str: "*"},
												},
											},
											Lexpr: &ast.A_Const{
												Val: &ast.Integer{Ival: 3},
											},
											Rexpr: &ast.A_Const{
												Val: &ast.Integer{Ival: 4},
											},
										},
									},
								},
							},
						},
						FromClause: &ast.List{},
					},
				},
			},
		},

		// Select with From Clause tests
		{
			stmt: `SELECT * FROM users`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.ColumnRef{
										Fields: &ast.List{
											Items: []ast.Node{
												&ast.A_Star{},
											},
										},
									},
								},
							},
						},
						FromClause: &ast.List{
							Items: []ast.Node{
								&ast.RangeVar{
									Relname: strPtr("users"),
								},
							},
						},
					},
				},
			},
		},
		{
			stmt: "SELECT id AS identifier FROM users",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Name: strPtr("identifier"),
									Val: &ast.ColumnRef{
										Fields: &ast.List{
											Items: []ast.Node{
												&ast.String{Str: "id"},
											},
										},
									},
								},
							},
						},
						FromClause: &ast.List{
							Items: []ast.Node{
								&ast.RangeVar{
									Relname: strPtr("users"),
								},
							},
						},
					},
				},
			},
		},
		{
			stmt: "SELECT a.b.c FROM table",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.ColumnRef{
										Fields: &ast.List{
											Items: []ast.Node{
												&ast.String{Str: "a"},
												&ast.String{Str: "b"},
												&ast.String{Str: "c"},
											},
										},
									},
								},
							},
						},
						FromClause: &ast.List{
							Items: []ast.Node{
								&ast.RangeVar{
									Relname: strPtr("table"),
								},
							},
						},
					},
				},
			},
		},
		{
			stmt: "SELECT id.age, 3.14, 'abc', NULL, false FROM users",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.ColumnRef{
										Fields: &ast.List{
											Items: []ast.Node{
												&ast.String{Str: "id"},
												&ast.String{Str: "age"},
											},
										},
									},
								},
								&ast.ResTarget{
									Val: &ast.A_Const{
										Val: &ast.Float{Str: "3.14"},
									},
								},
								&ast.ResTarget{
									Val: &ast.A_Const{
										Val: &ast.String{Str: "abc"},
									},
								},
								&ast.ResTarget{
									Val: &ast.Null{},
								},
								&ast.ResTarget{
									Val: &ast.A_Const{
										Val: &ast.Boolean{Boolval: false},
									},
								},
							},
						},
						FromClause: &ast.List{
							Items: []ast.Node{
								&ast.RangeVar{
									Relname: strPtr("users"),
								},
							},
						},
					},
				},
			},
		},
		{
			stmt: `SELECT id, name FROM users WHERE age > 30`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.ColumnRef{
										Fields: &ast.List{
											Items: []ast.Node{
												&ast.String{Str: "id"},
											},
										},
									},
								},
								&ast.ResTarget{
									Val: &ast.ColumnRef{
										Fields: &ast.List{
											Items: []ast.Node{
												&ast.String{Str: "name"},
											},
										},
									},
								},
							},
						},
						FromClause: &ast.List{
							Items: []ast.Node{
								&ast.RangeVar{
									Relname: strPtr("users"),
								},
							},
						},
						WhereClause: &ast.A_Expr{
							Name: &ast.List{
								Items: []ast.Node{
									&ast.String{Str: ">"},
								},
							},
							Lexpr: &ast.ColumnRef{
								Fields: &ast.List{
									Items: []ast.Node{
										&ast.String{Str: "age"},
									},
								},
							},
							Rexpr: &ast.A_Const{
								Val: &ast.Integer{Ival: 30},
							},
						},
					},
				},
			},
		},
	}

	p := ydb.NewParser()

	for _, tc := range tests {
		t.Run(tc.stmt, func(t *testing.T) {
			stmts, err := p.Parse(strings.NewReader(tc.stmt))
			if err != nil {
				t.Fatalf("Failed to parse query %q: %v", tc.stmt, err)
			}
			if len(stmts) == 0 {
				t.Fatalf("Query %q was not parsed", tc.stmt)
			}

			diff := cmp.Diff(tc.expected, &stmts[0],
				cmpopts.IgnoreFields(ast.RawStmt{}, "StmtLocation", "StmtLen"),
				// cmpopts.IgnoreFields(ast.SelectStmt{}, "Location"),
				cmpopts.IgnoreFields(ast.A_Const{}, "Location"),
				cmpopts.IgnoreFields(ast.ResTarget{}, "Location"),
				cmpopts.IgnoreFields(ast.ColumnRef{}, "Location"),
				cmpopts.IgnoreFields(ast.A_Expr{}, "Location"),
				cmpopts.IgnoreFields(ast.RangeVar{}, "Location"),
			)
			if diff != "" {
				t.Errorf("AST mismatch for %q (-expected +got):\n%s", tc.stmt, diff)
			}
		})
	}
}
