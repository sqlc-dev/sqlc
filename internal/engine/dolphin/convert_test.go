package dolphin

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
)

// walkStmts collects every node reachable from the parsed statements.
func walkStmts(t *testing.T, query string) []ast.Node {
	t.Helper()

	p := NewParser()
	stmts, err := p.Parse(strings.NewReader(query))
	if err != nil {
		t.Fatalf("parse %q: %v", query, err)
	}

	var nodes []ast.Node
	for i := range stmts {
		astutils.Walk(astutils.VisitorFunc(func(n ast.Node) {
			nodes = append(nodes, n)
		}), stmts[i].Raw.Stmt)
	}
	return nodes
}

func TestConvertUnaryExpr_Not(t *testing.T) {
	for _, tc := range []struct {
		name  string
		query string
	}{
		{
			name:  "NOT keyword",
			query: "SELECT personid FROM persons WHERE NOT sqlc.arg('foo')",
		},
		{
			name:  "bang operator",
			query: "SELECT personid FROM persons WHERE ! sqlc.arg('foo')",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			nodes := walkStmts(t, tc.query)

			var (
				foundNot bool
				foundArg bool
			)
			for _, n := range nodes {
				switch v := n.(type) {
				case *ast.TODO:
					t.Fatalf("query produced an opaque TODO node, sqlc.arg() would not be rewritten")
				case *ast.BoolExpr:
					if v.Boolop == ast.BoolExprTypeNot {
						foundNot = true
					}
				case *ast.FuncCall:
					if v.Func != nil && v.Func.Schema == "sqlc" && v.Func.Name == "arg" {
						foundArg = true
					}
				}
			}

			if !foundNot {
				t.Errorf("expected a BoolExpr with BoolExprTypeNot in the AST")
			}
			if !foundArg {
				t.Errorf("expected the nested sqlc.arg() FuncCall to be reachable in the AST")
			}
		})
	}
}
