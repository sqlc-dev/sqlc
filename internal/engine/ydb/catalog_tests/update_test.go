package ydb_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sqlc-dev/sqlc/internal/engine/ydb"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestUpdate(t *testing.T) {
	tests := []struct {
		stmt     string
		expected ast.Node
	}{
		{
			stmt: "UPDATE users SET name = 'Bob' WHERE id = 1 RETURNING id;",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.UpdateStmt{
						Relations: &ast.List{
							Items: []ast.Node{
								&ast.RangeVar{Relname: strPtr("users"), Inh: true},
							},
						},
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Name: strPtr("name"),
									Val: &ast.A_Const{
										Val: &ast.String{Str: "Bob"},
									},
								},
							},
						},
						WhereClause: &ast.A_Expr{
							Name: &ast.List{Items: []ast.Node{&ast.String{Str: "="}}},
							Lexpr: &ast.ColumnRef{
								Fields: &ast.List{Items: []ast.Node{&ast.String{Str: "id"}}},
							},
							Rexpr: &ast.A_Const{
								Val: &ast.Integer{Ival: 1},
							},
						},
						ReturningList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Indirection: &ast.List{},
									Val: &ast.ColumnRef{
										Fields: &ast.List{Items: []ast.Node{&ast.String{Str: "id"}}},
									},
								},
							},
						},
						FromClause:   &ast.List{},
						WithClause:   nil,
						Batch:        false,
						OnCols:       &ast.List{Items: []ast.Node{}},
						OnSelectStmt: nil,
					},
				},
			},
		},
		{
			stmt: "BATCH UPDATE users SET name = 'Charlie' WHERE id = 2 RETURNING *",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.UpdateStmt{
						Relations: &ast.List{
							Items: []ast.Node{
								&ast.RangeVar{Relname: strPtr("users"), Inh: true},
							},
						},
						TargetList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Name: strPtr("name"),
									Val:  &ast.A_Const{Val: &ast.String{Str: "Charlie"}},
								},
							},
						},
						WhereClause: &ast.A_Expr{
							Name: &ast.List{Items: []ast.Node{&ast.String{Str: "="}}},
							Lexpr: &ast.ColumnRef{
								Fields: &ast.List{Items: []ast.Node{&ast.String{Str: "id"}}},
							},
							Rexpr: &ast.A_Const{
								Val: &ast.Integer{Ival: 2},
							},
						},
						ReturningList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Indirection: &ast.List{},
									Val: &ast.ColumnRef{
										Fields: &ast.List{Items: []ast.Node{&ast.A_Star{}}},
									},
								},
							},
						},
						FromClause:   &ast.List{},
						WithClause:   nil,
						Batch:        true,
						OnCols:       &ast.List{Items: []ast.Node{}},
						OnSelectStmt: nil,
					},
				},
			},
		},
		{
			stmt: "UPDATE users ON (id) VALUES (5) RETURNING id",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.UpdateStmt{
						Relations: &ast.List{Items: []ast.Node{&ast.RangeVar{Relname: strPtr("users"), Inh: true}}},
						OnCols: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{Name: strPtr("id")},
							},
						},
						OnSelectStmt: &ast.SelectStmt{
							DistinctClause: &ast.List{},
							ValuesLists: &ast.List{
								Items: []ast.Node{
									&ast.List{
										Items: []ast.Node{
											&ast.A_Const{Val: &ast.Integer{Ival: 5}},
										},
									},
								},
							},
							FromClause:    &ast.List{},
							TargetList:    &ast.List{},
							GroupClause:   &ast.List{},
							WindowClause:  &ast.List{},
							SortClause:    &ast.List{},
							LockingClause: &ast.List{},
						},
						ReturningList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Indirection: &ast.List{},
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
						FromClause:  &ast.List{},
						WithClause:  nil,
						Batch:       false,
						TargetList:  &ast.List{Items: []ast.Node{}},
						WhereClause: nil,
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
			)
			if diff != "" {
				t.Errorf("AST mismatch for %q (-expected +got):\n%s", tc.stmt, diff)
			}
		})
	}
}
