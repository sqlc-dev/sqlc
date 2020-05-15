package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type SQLValueFunction struct {
	Xpr      ast.Node
	Op       SQLValueFunctionOp
	Type     Oid
	Typmod   int32
	Location int
}

func (n *SQLValueFunction) Pos() int {
	return n.Location
}
