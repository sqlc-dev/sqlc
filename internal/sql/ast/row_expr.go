package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type RowExpr struct {
	Xpr       ast.Node
	Args      *List
	RowTypeid Oid
	RowFormat CoercionForm
	Colnames  *List
	Location  int
}

func (n *RowExpr) Pos() int {
	return n.Location
}
