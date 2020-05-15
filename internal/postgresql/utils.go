package postgresql

import (
	nodes "github.com/lfittl/pg_query_go/nodes"

	"github.com/kyleconroy/sqlc/internal/postgresql/ast"
)

func isNotNull(n nodes.ColumnDef) bool {
	if n.IsNotNull {
		return true
	}
	for _, c := range n.Constraints.Items {
		switch n := c.(type) {
		case nodes.Constraint:
			if n.Contype == nodes.CONSTR_NOTNULL {
				return true
			}
			if n.Contype == nodes.CONSTR_PRIMARY {
				return true
			}
		}
	}
	return false
}

func IsNamedParamFunc(node nodes.Node) bool {
	fun, ok := node.(nodes.FuncCall)
	return ok && ast.Join(fun.Funcname, ".") == "sqlc.arg"
}

func IsNamedParamSign(node nodes.Node) bool {
	expr, ok := node.(nodes.A_Expr)
	return ok && ast.Join(expr.Name, ".") == "@"
}
