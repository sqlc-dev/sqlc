package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type FuncExpr struct {
	Xpr            ast.Node
	Funcid         Oid
	Funcresulttype Oid
	Funcretset     bool
	Funcvariadic   bool
	Funcformat     CoercionForm
	Funccollid     Oid
	Inputcollid    Oid
	Args           *ast.List
	Location       int
}

func (n *FuncExpr) Pos() int {
	return n.Location
}
