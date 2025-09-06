package ydb_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sqlc-dev/sqlc/internal/engine/ydb"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestPragma(t *testing.T) {
	tests := []struct {
		stmt     string
		expected ast.Node
	}{
		{
			stmt: `PRAGMA AutoCommit`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.Pragma_stmt{
						Name: &ast.List{
							Items: []ast.Node{
								&ast.A_Const{Val: &ast.String{Str: "autocommit"}},
							},
						},
					},
				},
			},
		},
		{
			stmt: `PRAGMA TablePathPrefix = "home/yql"`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.Pragma_stmt{
						Name: &ast.List{
							Items: []ast.Node{
								&ast.A_Const{Val: &ast.String{Str: "tablepathprefix"}},
							},
						},
						Equals: true,
						Values: &ast.List{
							Items: []ast.Node{
								&ast.A_Const{Val: &ast.String{Str: "home/yql"}},
							},
						},
					},
				},
			},
		},
		{
			stmt: `PRAGMA Warning("disable", "1101")`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.Pragma_stmt{
						Name: &ast.List{
							Items: []ast.Node{
								&ast.A_Const{Val: &ast.String{Str: "warning"}},
							},
						},
						Equals: false,
						Values: &ast.List{
							Items: []ast.Node{
								&ast.A_Const{Val: &ast.String{Str: "disable"}},
								&ast.A_Const{Val: &ast.String{Str: "1101"}},
							},
						},
					},
				},
			},
		},
		{
			stmt: `PRAGMA yson.AutoConvert = true`,
			expected: &ast.Statement{
				Raw: &ast.RawStmt{
					Stmt: &ast.Pragma_stmt{
						Name: &ast.List{
							Items: []ast.Node{
								&ast.A_Const{Val: &ast.String{Str: "yson"}},
								&ast.A_Const{Val: &ast.String{Str: "autoconvert"}},
							},
						},
						Equals: true,
						Values: &ast.List{
							Items: []ast.Node{
								&ast.A_Const{Val: &ast.Boolean{Boolval: true}},
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
				cmpopts.IgnoreFields(ast.Pragma_stmt{}, "Location"),
				cmpopts.IgnoreFields(ast.ColumnRef{}, "Location"),
				cmpopts.IgnoreFields(ast.A_Const{}, "Location"),
			)
			if diff != "" {
				t.Errorf("AST mismatch for %q (-expected +got):\n%s", tc.stmt, diff)
			}
		})
	}
}
