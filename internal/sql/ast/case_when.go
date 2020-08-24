package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CaseWhen struct {
	Xpr      Node
	Expr     Node
	Result   Node
	Location int
}

func (n *CaseWhen) Pos() int {
	return n.Location
}
