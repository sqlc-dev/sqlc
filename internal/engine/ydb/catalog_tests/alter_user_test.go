package ydb_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sqlc-dev/sqlc/internal/engine/ydb"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestAlterUser(t *testing.T) {
	tests := []struct {
		stmt     string
		expected ast.Node
	}{
		{
			stmt: `ALTER USER alice RENAME TO queen`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.AlterRoleStmt{
						Role: &ast.RoleSpec{
							Rolename: strPtr("alice"),
							Roletype: ast.RoleSpecType(1),
						},
						Action: 1,
						Options: &ast.List{
							Items: []ast.Node{
								&ast.DefElem{
									Defname:   strPtr("rename"),
									Arg:       &ast.String{Str: "queen"},
									Defaction: ast.DefElemAction(1),
								},
							},
						},
					},
				},
			},
		},
		{
			stmt: `ALTER USER bob LOGIN`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.AlterRoleStmt{
						Role: &ast.RoleSpec{
							Rolename: strPtr("bob"),
							Roletype: ast.RoleSpecType(1),
						},
						Action: 1,
						Options: &ast.List{
							Items: []ast.Node{
								&ast.DefElem{
									Defname: strPtr("login"),
									Arg:     &ast.Boolean{Boolval: true},
								},
							},
						},
					},
				},
			},
		},
		{
			stmt: `ALTER USER charlie NOLOGIN`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.AlterRoleStmt{
						Role: &ast.RoleSpec{
							Rolename: strPtr("charlie"),
							Roletype: ast.RoleSpecType(1),
						},
						Action: 1,
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
			stmt: `ALTER USER dave PASSWORD 'qwerty'`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.AlterRoleStmt{
						Role: &ast.RoleSpec{
							Rolename: strPtr("dave"),
							Roletype: ast.RoleSpecType(1),
						},
						Action: 1,
						Options: &ast.List{
							Items: []ast.Node{
								&ast.DefElem{
									Defname: strPtr("password"),
									Arg:     &ast.String{Str: "qwerty"},
								},
							},
						},
					},
				},
			},
		},
		{
			stmt: `ALTER USER elena HASH 'abc123'`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.AlterRoleStmt{
						Role: &ast.RoleSpec{
							Rolename: strPtr("elena"),
							Roletype: ast.RoleSpecType(1),
						},
						Action: 1,
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
				t.Fatalf("Failed to parse query %q: %v", tc.stmt, err)
			}
			if len(stmts) == 0 {
				t.Fatalf("Query %q was not parsed", tc.stmt)
			}

			diff := cmp.Diff(tc.expected, &stmts[0],
				cmpopts.IgnoreFields(ast.RawStmt{}, "StmtLocation", "StmtLen"),
				cmpopts.IgnoreFields(ast.DefElem{}, "Location"),
				cmpopts.IgnoreFields(ast.RoleSpec{}, "Location"),
				cmpopts.IgnoreFields(ast.A_Const{}, "Location"),
			)
			if diff != "" {
				t.Errorf("AST mismatch for %q (-expected +got):\n%s", tc.stmt, diff)
			}
		})
	}
}
