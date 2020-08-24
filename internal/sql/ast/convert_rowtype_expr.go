package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type ConvertRowtypeExpr struct {
	Xpr           ast.Node
	Arg           ast.Node
	Resulttype    Oid
	Convertformat CoercionForm
	Location      int
}

func (n *ConvertRowtypeExpr) Pos() int {
	return n.Location
}
