package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type BoolExpr struct {
	Xpr      ast.Node
	Boolop   BoolExprType
	Args     *List
	Location int
}

func (n *BoolExpr) Pos() int {
	return n.Location
}
