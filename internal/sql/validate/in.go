package validate

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

type inVisitor struct {
	catalog *catalog.Catalog
	err     error
}

func (v *inVisitor) Visit(node ast.Node) astutils.Visitor {
	if v.err != nil {
		return nil
	}

	in, ok := node.(*ast.In)
	if !ok {
		return v
	}

	// Validate that sqlc.slice in an IN statement is the only arg, eg:
	//    id IN (sqlc.slice("ids"))       -- GOOD
	//    id in (0, 1, sqlc.slice("ids")) -- BAD

	if len(in.List) <= 1 {
		return v
	}

	for _, n := range in.List {
		call, ok := n.(*ast.FuncCall)
		if !ok {
			continue
		}
		fn := call.Func
		if fn == nil {
			continue
		}

		if fn.Schema == "sqlc" && fn.Name == "slice" {
			var inExpr, sliceArg string

			// determine inExpr
			switch n := in.Expr.(type) {
			case *ast.ColumnRef:
				inExpr = n.Name
			default:
				inExpr = "..."
			}

			// determine sliceArg
			if len(call.Args.Items) == 1 {
				switch n := call.Args.Items[0].(type) {
				case *ast.A_Const:
					if str, ok := n.Val.(*ast.String); ok {
						sliceArg = "\"" + str.Str + "\""
					} else {
						sliceArg = "?"
					}
				case *ast.ColumnRef:
					sliceArg = n.Name
				default:
					// impossible, validate.FuncCall should have caught this
					sliceArg = "..."
				}
			}
			v.err = &sqlerr.Error{
				Message:  fmt.Sprintf("expected '%s IN' expr to consist only of sqlc.slice(%s); eg ", inExpr, sliceArg),
				Location: call.Pos(),
			}
		}
	}

	return v
}

func In(c *catalog.Catalog, n ast.Node) error {
	visitor := inVisitor{catalog: c}
	astutils.Walk(&visitor, n)
	return visitor.err
}
