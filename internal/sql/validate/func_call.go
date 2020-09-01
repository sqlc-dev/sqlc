package validate

import (
	"errors"
	"fmt"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"
)

type funcCallVisitor struct {
	catalog *catalog.Catalog
	err     error
}

func (v *funcCallVisitor) Visit(node ast.Node) astutils.Visitor {
	if v.err != nil {
		return nil
	}

	call, ok := node.(*ast.FuncCall)
	if !ok {
		return v
	}
	fn := call.Func
	if fn == nil {
		return v
	}

	// Custom validation for sqlc.arg
	// TODO: Replace this once type-checking is implemented
	if fn.Schema == "sqlc" {
		if fn.Name != "arg" {
			v.err = sqlerr.FunctionNotFound("sqlc." + fn.Name)
			return nil
		}
		if call.Args == nil || len(call.Args.Items) == 0 {
			return v
		}
		if len(call.Args.Items) > 1 {
			v.err = &sqlerr.Error{
				Message:  fmt.Sprintf("expected 1 parameter to sqlc.arg; got %d", len(call.Args.Items)),
				Location: call.Pos(),
			}
			return nil
		}
		switch n := call.Args.Items[0].(type) {
		case *ast.A_Const:
		case *ast.ColumnRef:
		default:
			v.err = &sqlerr.Error{
				Message:  fmt.Sprintf("expected parameter to sqlc.arg to be string or reference; got %T", n),
				Location: call.Pos(),
			}
			return nil
		}
	}

	fun, err := v.catalog.ResolveFuncCall(call)
	if fun != nil || errors.Is(err, sqlerr.NotFound) {
		return v
	}

	v.err = err
	return nil
}

func FuncCall(c *catalog.Catalog, n ast.Node) error {
	visitor := funcCallVisitor{catalog: c}
	astutils.Walk(&visitor, n)
	return visitor.err
}
