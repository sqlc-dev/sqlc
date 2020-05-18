package named

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/ast/pg"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
)

func IsParamFunc(node ast.Node) bool {
	fun, ok := node.(*pg.FuncCall)
	return ok && astutils.Join(fun.Funcname, ".") == "sqlc.arg"
}

func IsParamSign(node ast.Node) bool {
	expr, ok := node.(*pg.A_Expr)
	return ok && astutils.Join(expr.Name, ".") == "@"
}
