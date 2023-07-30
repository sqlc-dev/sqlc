package named

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
)

// IsParamFunc fulfills the astutils.Search
func IsParamFunc(node ast.Node) bool {
	call, ok := node.(*ast.FuncCall)
	if !ok {
		return false
	}

	if call.Func == nil {
		return false
	}

	isValid := call.Func.Schema == "sqlc" && (call.Func.Name == "arg" || call.Func.Name == "narg" || call.Func.Name == "slice")
	return isValid
}

func IsParamSign(node ast.Node) bool {
	expr, ok := node.(*ast.A_Expr)
	return ok && astutils.Join(expr.Name, ".") == "@"
}
