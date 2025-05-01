package postgresql

import (
	nodes "github.com/pganalyze/pg_query_go/v6"
)

func isArray(n *nodes.TypeName) bool {
	if n == nil {
		return false
	}
	return len(n.ArrayBounds) > 0
}

func isNotNull(n *nodes.ColumnDef) bool {
	if n.IsNotNull {
		return true
	}
	for _, c := range n.Constraints {
		switch inner := c.Node.(type) {
		case *nodes.Node_Constraint:
			if inner.Constraint.Contype == nodes.ConstrType_CONSTR_NOTNULL {
				return true
			}
			if inner.Constraint.Contype == nodes.ConstrType_CONSTR_PRIMARY {
				return true
			}
		}
	}
	return false
}

func IsNamedParamFunc(node *nodes.Node) bool {
	fun, ok := node.Node.(*nodes.Node_FuncCall)
	return ok && joinNodes(fun.FuncCall.Funcname, ".") == "sqlc.arg"
}

func IsNamedParamSign(node *nodes.Node) bool {
	expr, ok := node.Node.(*nodes.Node_AExpr)
	return ok && joinNodes(expr.AExpr.Name, ".") == "@"
}

func makeByte(s string) byte {
	var b byte
	if s == "" {
		return b
	}
	return []byte(s)[0]
}

func makeUint32Slice(in []uint64) []uint32 {
	out := make([]uint32, len(in))
	for i, v := range in {
		out[i] = uint32(v)
	}
	return out
}

func makeString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
