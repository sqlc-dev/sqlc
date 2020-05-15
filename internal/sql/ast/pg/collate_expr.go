package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CollateExpr struct {
	Xpr      ast.Node
	Arg      ast.Node
	CollOid  Oid
	Location int
}

func (n *CollateExpr) Pos() int {
	return n.Location
}
