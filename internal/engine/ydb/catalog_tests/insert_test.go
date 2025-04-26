package ydb_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sqlc-dev/sqlc/internal/engine/ydb"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestInsert(t *testing.T) {
	tests := []struct {
		stmt     string
		expected ast.Node
	}{
		{
			stmt: "INSERT INTO users (id, name) VALUES (1, 'Alice') RETURNING *",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.InsertStmt{
						Relation: &ast.RangeVar{Relname: strPtr("users")},
						Cols: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{Name: strPtr("id")},
								&ast.ResTarget{Name: strPtr("name")},
							},
						},
						SelectStmt: &ast.SelectStmt{
							ValuesLists: &ast.List{
								Items: []ast.Node{
									&ast.List{
										Items: []ast.Node{
											&ast.A_Const{Val: &ast.Integer{Ival: 1}},
											&ast.A_Const{Val: &ast.String{Str: "Alice"}},
										},
									},
								},
							},
						},
						OnConflictClause: &ast.OnConflictClause{},
						ReturningList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.ColumnRef{
										Fields: &ast.List{Items: []ast.Node{&ast.A_Star{}}},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			stmt: "INSERT OR IGNORE INTO users (id) VALUES (3) RETURNING id, name",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.InsertStmt{
						Relation: &ast.RangeVar{Relname: strPtr("users")},
						Cols: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{Name: strPtr("id")},
							},
						},
						SelectStmt: &ast.SelectStmt{
							ValuesLists: &ast.List{
								Items: []ast.Node{
									&ast.List{
										Items: []ast.Node{
											&ast.A_Const{Val: &ast.Integer{Ival: 3}},
										},
									},
								},
							},
						},
						OnConflictClause: &ast.OnConflictClause{
							Action: ast.OnConflictAction_INSERT_OR_IGNORE,
						},
						ReturningList: &ast.List{
							Items: []ast.Node{
								&ast.ResTarget{
									Val: &ast.ColumnRef{Fields: &ast.List{Items: []ast.Node{&ast.String{Str: "id"}}}},
								},
								&ast.ResTarget{
									Val: &ast.ColumnRef{Fields: &ast.List{Items: []ast.Node{&ast.String{Str: "name"}}}},
								},
							},
						},
					},
				},
			},
		},
		{
			stmt: "UPSERT INTO users (id) VALUES (4) RETURNING id",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.InsertStmt{
						Relation:         &ast.RangeVar{Relname: strPtr("users")},
						Cols:             &ast.List{Items: []ast.Node{&ast.ResTarget{Name: strPtr("id")}}},
						SelectStmt:       &ast.SelectStmt{ValuesLists: &ast.List{Items: []ast.Node{&ast.List{Items: []ast.Node{&ast.A_Const{Val: &ast.Integer{Ival: 4}}}}}}},
						OnConflictClause: &ast.OnConflictClause{Action: ast.OnConflictAction_UPSERT},
						ReturningList:    &ast.List{Items: []ast.Node{&ast.ResTarget{Val: &ast.ColumnRef{Fields: &ast.List{Items: []ast.Node{&ast.String{Str: "id"}}}}}}},
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
				t.Fatalf("Ошибка парсинга запроса %q: %v", tc.stmt, err)
			}
			if len(stmts) == 0 {
				t.Fatalf("Запрос %q не распарсен", tc.stmt)
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
				t.Errorf("Несовпадение AST (-ожидалось +получено):\n%s", diff)
			}
		})
	}
}
