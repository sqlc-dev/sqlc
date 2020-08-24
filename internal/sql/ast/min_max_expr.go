package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type MinMaxExpr struct {
	Xpr          ast.Node
	Minmaxtype   Oid
	Minmaxcollid Oid
	Inputcollid  Oid
	Op           MinMaxOp
	Args         *List
	Location     int
}

func (n *MinMaxExpr) Pos() int {
	return n.Location
}
