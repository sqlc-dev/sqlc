package named

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
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

	// sqlite doesn't support the sql.narg syntax and the parser fails, so we have to "sqlc_narg"
	isValid := (call.Func.Schema == "sqlc" && (call.Func.Name == "arg" || call.Func.Name == "narg")) || call.Func.Name == "sqlc_narg"
	return isValid
}

func IsParamSign(node ast.Node) bool {
	expr, ok := node.(*ast.A_Expr)
	return ok && astutils.Join(expr.Name, ".") == "@"
}
