package ydb_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sqlc-dev/sqlc/internal/engine/ydb"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		stmt     string
		expected ast.Node
	}{
		{
			stmt: `CREATE USER alice`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.CreateRoleStmt{
						StmtType: ast.RoleStmtType(2), // CREATE USER
						Role:     strPtr("alice"),
						Options:  &ast.List{},
					},
				},
			},
		},
		{
			stmt: `CREATE USER bob PASSWORD 'secret'`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.CreateRoleStmt{
						StmtType: ast.RoleStmtType(2),
						Role:     strPtr("bob"),
						Options: &ast.List{
							Items: []ast.Node{
								&ast.DefElem{
									Defname: strPtr("password"),
									Arg:     &ast.String{Str: "secret"},
								},
							},
						},
					},
				},
			},
		},
		{
			stmt: `CREATE USER charlie LOGIN`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.CreateRoleStmt{
						StmtType: ast.RoleStmtType(2),
						Role:     strPtr("charlie"),
						Options: &ast.List{
							Items: []ast.Node{
								&ast.DefElem{
									Defname: strPtr("login"),
									Arg:      &ast.Boolean{Boolval: true},
								},
							},
						},
					},
				},
			},
		},
		{
			stmt: `CREATE USER dave NOLOGIN`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.CreateRoleStmt{
						StmtType: ast.RoleStmtType(2),
						Role:     strPtr("dave"),
						Options: &ast.List{
							Items: []ast.Node{
								&ast.DefElem{
									Defname: strPtr("nologin"),
									Arg:     &ast.Boolean{Boolval: false},
								},
							},
						},
					},
				},
			},
		},
		{
			stmt: `CREATE USER bjorn HASH 'abc123'`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.CreateRoleStmt{
						StmtType: ast.RoleStmtType(2),
						Role:     strPtr("bjorn"),
						Options: &ast.List{
							Items: []ast.Node{
								&ast.DefElem{
									Defname: strPtr("hash"),
									Arg:     &ast.String{Str: "abc123"},
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
				cmpopts.IgnoreFields(ast.A_Const{}, "Location"),
				cmpopts.IgnoreFields(ast.DefElem{}, "Location"),
			)
			if diff != "" {
				t.Errorf("Несовпадение AST (-ожидалось +получено):\n%s", diff)
			}
		})
	}
}
