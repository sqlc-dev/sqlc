package ydb_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sqlc-dev/sqlc/internal/engine/ydb"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestCreateGroup(t *testing.T) {
	tests := []struct {
		stmt     string
		expected ast.Node
	}{
		{
			stmt: `CREATE GROUP group1`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.CreateRoleStmt{
						StmtType: ast.RoleStmtType(3), // CREATE GROUP
						Role:     strPtr("group1"),
						Options:  &ast.List{},
					},
				},
			},
		},
		{
			stmt: `CREATE GROUP group1 WITH USER alice`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.CreateRoleStmt{
						StmtType: ast.RoleStmtType(3),
						Role:     strPtr("group1"),
						Options: &ast.List{
							Items: []ast.Node{
								&ast.DefElem{
									Defname: strPtr("rolemembers"),
									Arg: &ast.List{
										Items: []ast.Node{
											&ast.RoleSpec{
												Roletype: ast.RoleSpecType(1),
												Rolename: strPtr("alice"),
											},
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
			stmt: `CREATE GROUP group1 WITH USER alice, bebik`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.CreateRoleStmt{
						StmtType: ast.RoleStmtType(3),
						Role:     strPtr("group1"),
						Options: &ast.List{
							Items: []ast.Node{
								&ast.DefElem{
									Defname: strPtr("rolemembers"),
									Arg: &ast.List{
										Items: []ast.Node{
											&ast.RoleSpec{
												Roletype: ast.RoleSpecType(1),
												Rolename: strPtr("alice"),
											},
											&ast.RoleSpec{
												Roletype: ast.RoleSpecType(1),
												Rolename: strPtr("bebik"),
											},
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
