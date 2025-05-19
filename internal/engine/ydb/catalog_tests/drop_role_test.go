package ydb_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sqlc-dev/sqlc/internal/engine/ydb"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestDropRole(t *testing.T) {
	tests := []struct {
		stmt     string
		expected ast.Node
	}{
		{
			stmt: `DROP USER user1;`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.DropRoleStmt{
						MissingOk: false,
						Roles: &ast.List{
							Items: []ast.Node{
								&ast.RoleSpec{Rolename: strPtr("user1"), Roletype: ast.RoleSpecType(1)},
							},
						},
					},
				},
			},
		},
		{
			stmt: "DROP USER IF EXISTS admin, user2",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.DropRoleStmt{
						MissingOk: true,
						Roles: &ast.List{
							Items: []ast.Node{
								&ast.RoleSpec{Rolename: strPtr("admin"), Roletype: ast.RoleSpecType(1)},
								&ast.RoleSpec{Rolename: strPtr("user2"), Roletype: ast.RoleSpecType(1)},
							},
						},
					},
				},
			},
		},
		{
			stmt: "DROP GROUP team1, team2",
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.DropRoleStmt{
						MissingOk: false,
						Roles: &ast.List{
							Items: []ast.Node{
								&ast.RoleSpec{Rolename: strPtr("team1"), Roletype: ast.RoleSpecType(1)},
								&ast.RoleSpec{Rolename: strPtr("team2"), Roletype: ast.RoleSpecType(1)},
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
				t.Fatalf("Error parsing %q: %v", tc.stmt, err)
			}
			if len(stmts) == 0 {
				t.Fatalf("Statement %q was not parsed", tc.stmt)
			}

			diff := cmp.Diff(tc.expected, &stmts[0],
				cmpopts.IgnoreFields(ast.RawStmt{}, "StmtLocation", "StmtLen"),
				cmpopts.IgnoreFields(ast.RoleSpec{}, "Location"),
			)
			if diff != "" {
				t.Errorf("AST mismatch for %q (-expected +got):\n%s", tc.stmt, diff)
			}
		})
	}
}
