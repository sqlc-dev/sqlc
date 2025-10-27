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
						DistinctClause: &ast.List{},
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.A_Const{
										Val: &ast.Integer{Ival: 52},
									},
								},
							},
						},
						FromClause:    &ast.List{},
						GroupClause:   &ast.List{},
						WindowClause:  &ast.List{},
						ValuesLists:   &ast.List{},
						SortClause:    &ast.List{},
						LockingClause: &ast.List{},
					},
				},
			},
		},
		{
			stmt: `SELECT 'hello'`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						DistinctClause: &ast.List{},
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.A_Const{
										Val: &ast.String{Str: "hello"},
									},
								},
							},
						},
						FromClause:    &ast.List{},
						GroupClause:   &ast.List{},
						WindowClause:  &ast.List{},
						ValuesLists:   &ast.List{},
						SortClause:    &ast.List{},
						LockingClause: &ast.List{},
					},
				},
			},
		},
		{
			stmt: `SELECT 'it\'s string with quote in it'`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						DistinctClause: &ast.List{},
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.A_Const{
										Val: &ast.String{Str: `it\'s string with quote in it`},
									},
								},
							},
						},
						FromClause:    &ast.List{},
						GroupClause:   &ast.List{},
						WindowClause:  &ast.List{},
						ValuesLists:   &ast.List{},
						SortClause:    &ast.List{},
						LockingClause: &ast.List{},
					},
				},
			},
		},
		{
			stmt: "SELECT 3.14",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						DistinctClause: &ast.List{},
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.A_Const{
										Val: &ast.Float{Str: "3.14"},
									},
								},
							},
						},
						FromClause:    &ast.List{},
						GroupClause:   &ast.List{},
						WindowClause:  &ast.List{},
						ValuesLists:   &ast.List{},
						SortClause:    &ast.List{},
						LockingClause: &ast.List{},
					},
				},
			},
		},
		{
			stmt: "SELECT NULL",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						DistinctClause: &ast.List{},
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.Null{},
								},
							},
						},
						FromClause:    &ast.List{},
						GroupClause:   &ast.List{},
						WindowClause:  &ast.List{},
						ValuesLists:   &ast.List{},
						SortClause:    &ast.List{},
						LockingClause: &ast.List{},
					},
				},
			},
		},
		{
			stmt: "SELECT true",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						DistinctClause: &ast.List{},
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.A_Const{
										Val: &ast.Boolean{Boolval: true},
									},
								},
							},
						},
						FromClause:    &ast.List{},
						GroupClause:   &ast.List{},
						WindowClause:  &ast.List{},
						ValuesLists:   &ast.List{},
						SortClause:    &ast.List{},
						LockingClause: &ast.List{},
					},
				},
			},
		},
		{
			stmt: "SELECT 2+3*4",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						DistinctClause: &ast.List{},
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
						FromClause:    &ast.List{},
						GroupClause:   &ast.List{},
						WindowClause:  &ast.List{},
						ValuesLists:   &ast.List{},
						SortClause:    &ast.List{},
						LockingClause: &ast.List{},
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
						DistinctClause: &ast.List{},
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
									Inh:     true,
								},
							},
						},
						GroupClause:   &ast.List{},
						WindowClause:  &ast.List{},
						ValuesLists:   &ast.List{},
						SortClause:    &ast.List{},
						LockingClause: &ast.List{},
					},
				},
			},
		},
		{
			stmt: "SELECT id AS identifier FROM users",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						DistinctClause: &ast.List{},
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
									Inh:     true,
								},
							},
						},
						GroupClause:   &ast.List{},
						WindowClause:  &ast.List{},
						ValuesLists:   &ast.List{},
						SortClause:    &ast.List{},
						LockingClause: &ast.List{},
					},
				},
			},
		},
		{
			stmt: "SELECT a.b.c FROM table",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						DistinctClause: &ast.List{},
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
									Inh:     true,
								},
							},
						},
						GroupClause:   &ast.List{},
						WindowClause:  &ast.List{},
						ValuesLists:   &ast.List{},
						SortClause:    &ast.List{},
						LockingClause: &ast.List{},
					},
				},
			},
		},
		{
			stmt: "SELECT id.age, 3.14, 'abc', NULL, false FROM users",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						DistinctClause: &ast.List{},
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
									Inh:     true,
								},
							},
						},
						GroupClause:   &ast.List{},
						WindowClause:  &ast.List{},
						ValuesLists:   &ast.List{},
						SortClause:    &ast.List{},
						LockingClause: &ast.List{},
					},
				},
			},
		},
		{
			stmt: `SELECT id, name FROM users WHERE age > 30`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						DistinctClause: &ast.List{},
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
									Inh:     true,
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
						GroupClause:   &ast.List{},
						WindowClause:  &ast.List{},
						ValuesLists:   &ast.List{},
						SortClause:    &ast.List{},
						LockingClause: &ast.List{},
					},
				},
			},
		},
		{
			stmt: `(SELECT 1) UNION ALL (SELECT 2)`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						DistinctClause: &ast.List{},
						TargetList:     &ast.List{},
						FromClause:     &ast.List{},
						GroupClause:    &ast.List{},
						WindowClause:   &ast.List{},
						ValuesLists:    &ast.List{},
						SortClause:     &ast.List{},
						LockingClause:  &ast.List{},
						Op:             ast.Union,
						All:            true,
						Larg: &ast.SelectStmt{
							DistinctClause: &ast.List{},
							TargetList: &ast.List{
								Items: []ast.Node{
									&ast.ResTarget{
										Val: &ast.A_Const{
											Val: &ast.Integer{Ival: 1},
										},
									},
								},
							},
							FromClause:    &ast.List{},
							GroupClause:   &ast.List{},
							WindowClause:  &ast.List{},
							ValuesLists:   &ast.List{},
							SortClause:    &ast.List{},
							LockingClause: &ast.List{},
						},
						Rarg: &ast.SelectStmt{
							DistinctClause: &ast.List{},
							TargetList: &ast.List{
								Items: []ast.Node{
									&ast.ResTarget{
										Val: &ast.A_Const{
											Val: &ast.Integer{Ival: 2},
										},
									},
								},
							},
							FromClause:    &ast.List{},
							GroupClause:   &ast.List{},
							WindowClause:  &ast.List{},
							ValuesLists:   &ast.List{},
							SortClause:    &ast.List{},
							LockingClause: &ast.List{},
						},
					},
				},
			},
		},
		{
			stmt: `SELECT id FROM users ORDER BY id DESC`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						DistinctClause: &ast.List{},
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
							},
						},
						FromClause: &ast.List{
							Items: []ast.Node{
								&ast.RangeVar{
									Relname: strPtr("users"),
									Inh:     true,
								},
							},
						},
						GroupClause:   &ast.List{},
						WindowClause:  &ast.List{},
						ValuesLists:   &ast.List{},
						LockingClause: &ast.List{},
						SortClause: &ast.List{
							Items: []ast.Node{
								&ast.SortBy{
									Node: &ast.ColumnRef{
										Fields: &ast.List{
											Items: []ast.Node{
												&ast.String{Str: "id"},
											},
										},
									},
									SortbyDir: ast.SortByDirDesc,
									UseOp:     &ast.List{},
								},
							},
						},
					},
				},
			},
		},
		{
			stmt: `SELECT id FROM users LIMIT 10 OFFSET 5`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						DistinctClause: &ast.List{},
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
							},
						},
						FromClause: &ast.List{
							Items: []ast.Node{
								&ast.RangeVar{
									Relname: strPtr("users"),
									Inh:     true,
								},
							},
						},
						GroupClause:   &ast.List{},
						WindowClause:  &ast.List{},
						ValuesLists:   &ast.List{},
						SortClause:    &ast.List{},
						LockingClause: &ast.List{},
						LimitCount: &ast.A_Const{
							Val: &ast.Integer{Ival: 10},
						},
						LimitOffset: &ast.A_Const{
							Val: &ast.Integer{Ival: 5},
						},
					},
				},
			},
		},
		{
			stmt: `SELECT id FROM users WHERE id > 10 GROUP BY id HAVING id > 10`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						DistinctClause: &ast.List{},
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
							},
						},
						FromClause: &ast.List{
							Items: []ast.Node{
								&ast.RangeVar{
									Relname: strPtr("users"),
									Inh:     true,
								},
							},
						},
						GroupClause: &ast.List{
							Items: []ast.Node{
								&ast.ColumnRef{
									Fields: &ast.List{
										Items: []ast.Node{
											&ast.String{Str: "id"},
										},
									},
								},
							},
						},
						WindowClause:  &ast.List{},
						ValuesLists:   &ast.List{},
						SortClause:    &ast.List{},
						LockingClause: &ast.List{},
						WhereClause: &ast.A_Expr{
							Name: &ast.List{
								Items: []ast.Node{
									&ast.String{Str: ">"},
								},
							},
							Lexpr: &ast.ColumnRef{
								Fields: &ast.List{
									Items: []ast.Node{
										&ast.String{Str: "id"},
									},
								},
							},
							Rexpr: &ast.A_Const{
								Val: &ast.Integer{Ival: 10},
							},
						},
						HavingClause: &ast.A_Expr{
							Name: &ast.List{
								Items: []ast.Node{
									&ast.String{Str: ">"},
								},
							},
							Lexpr: &ast.ColumnRef{
								Fields: &ast.List{
									Items: []ast.Node{
										&ast.String{Str: "id"},
									},
								},
							},
							Rexpr: &ast.A_Const{
								Val: &ast.Integer{Ival: 10},
							},
						},
					},
				},
			},
		},
		{
			stmt: `SELECT id FROM users GROUP BY ROLLUP (id)`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.SelectStmt{
						DistinctClause: &ast.List{},
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
							},
						},
						FromClause: &ast.List{
							Items: []ast.Node{
								&ast.RangeVar{
									Relname: strPtr("users"),
									Inh:     true,
								},
							},
						},
						GroupClause: &ast.List{
							Items: []ast.Node{
								&ast.GroupingSet{
									Kind: 1, // T_GroupingSet: ROLLUP
									Content: &ast.List{
										Items: []ast.Node{
											&ast.ColumnRef{
												Fields: &ast.List{
													Items: []ast.Node{
														&ast.String{Str: "id"},
													},
												},
											},
										},
									},
								},
							},
						},
						WindowClause:  &ast.List{},
						ValuesLists:   &ast.List{},
						SortClause:    &ast.List{},
						LockingClause: &ast.List{},
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
				cmpopts.IgnoreFields(ast.A_Const{}, "Location"),
				cmpopts.IgnoreFields(ast.ResTarget{}, "Location"),
				cmpopts.IgnoreFields(ast.ColumnRef{}, "Location"),
				cmpopts.IgnoreFields(ast.A_Expr{}, "Location"),
				cmpopts.IgnoreFields(ast.RangeVar{}, "Location"),
				cmpopts.IgnoreFields(ast.SortBy{}, "Location"),
			)
			if diff != "" {
				t.Errorf("AST mismatch for %q (-expected +got):\n%s", tc.stmt, diff)
			}
		})
	}
}
