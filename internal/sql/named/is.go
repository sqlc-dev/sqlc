package named

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
)

func IsParamFunc(node ast.Node) bool {
	call, ok := node.(*ast.FuncCall)
	if !ok {
		return false
	}
	if call.Func == nil {
		return false
	}
	return call.Func.Schema == "sqlc" && call.Func.Name == "arg"
}

func IsParamSign(node ast.Node) bool {
	expr, ok := node.(*ast.A_Expr)
	return ok && astutils.Join(expr.Name, ".") == "@"
}
