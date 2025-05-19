package ydb_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sqlc-dev/sqlc/internal/engine/ydb"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestAlterGroup(t *testing.T) {
	tests := []struct {
		stmt     string
		expected ast.Node
	}{
		{
			stmt: `ALTER GROUP admins RENAME TO superusers`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.AlterRoleStmt{
						Role: &ast.RoleSpec{
							Rolename: strPtr("admins"),
							Roletype: ast.RoleSpecType(1),
						},
						Action: 1,
						Options: &ast.List{
							Items: []ast.Node{
								&ast.DefElem{
									Defname:   strPtr("rename"),
									Defaction: ast.DefElemAction(1),
									Arg:       &ast.String{Str: "superusers"},
								},
							},
						},
					},
				},
			},
		},
		{
			stmt: `ALTER GROUP devs ADD USER alice, bob, carol`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.AlterRoleStmt{
						Role: &ast.RoleSpec{
							Rolename: strPtr("devs"),
							Roletype: ast.RoleSpecType(1),
						},
						Action: 1,
						Options: &ast.List{
							Items: []ast.Node{
								&ast.DefElem{
									Defname:   strPtr("rolemembers"),
									Defaction: ast.DefElemAction(3),
									Arg: &ast.List{
										Items: []ast.Node{
											&ast.RoleSpec{Rolename: strPtr("alice"), Roletype: ast.RoleSpecType(1)},
											&ast.RoleSpec{Rolename: strPtr("bob"), Roletype: ast.RoleSpecType(1)},
											&ast.RoleSpec{Rolename: strPtr("carol"), Roletype: ast.RoleSpecType(1)},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			stmt: `ALTER GROUP ops DROP USER dan, erin`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.AlterRoleStmt{
						Role: &ast.RoleSpec{
							Rolename: strPtr("ops"),
							Roletype: ast.RoleSpecType(1),
						},
						Action: 1,
						Options: &ast.List{
							Items: []ast.Node{
								&ast.DefElem{
									Defname:   strPtr("rolemembers"),
									Defaction: ast.DefElemAction(4),
									Arg: &ast.List{
										Items: []ast.Node{
											&ast.RoleSpec{Rolename: strPtr("dan"), Roletype: ast.RoleSpecType(1)},
											&ast.RoleSpec{Rolename: strPtr("erin"), Roletype: ast.RoleSpecType(1)},
										},
									},
								},
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
				t.Fatalf("Ошибка парсинга запроса %q: %v", tc.stmt, err)
			}
			if len(stmts) == 0 {
				t.Fatalf("Запрос %q не распарсен", tc.stmt)
			}

			diff := cmp.Diff(tc.expected, &stmts[0],
				cmpopts.IgnoreFields(ast.RawStmt{}, "StmtLocation", "StmtLen"),
				cmpopts.IgnoreFields(ast.DefElem{}, "Location"),
				cmpopts.IgnoreFields(ast.RoleSpec{}, "Location"),
				cmpopts.IgnoreFields(ast.A_Const{}, "Location"),
			)
			if diff != "" {
				t.Errorf("Несовпадение AST (-ожидалось +получено):\n%s", diff)
			}
		})
	}
}
