package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type A_ArrayExpr struct {
	Elements *List
	Location int
}

func (n *A_ArrayExpr) Pos() int {
	return n.Location
}
