package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type RowExpr struct {
	Xpr       ast.Node
	Args      *ast.List
	RowTypeid Oid
	RowFormat CoercionForm
	Colnames  *ast.List
	Location  int
}

func (n *RowExpr) Pos() int {
	return n.Location
}
