package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type OpExpr struct {
	Xpr          ast.Node
	Opno         Oid
	Opfuncid     Oid
	Opresulttype Oid
	Opretset     bool
	Opcollid     Oid
	Inputcollid  Oid
	Args         *ast.List
	Location     int
}

func (n *OpExpr) Pos() int {
	return n.Location
}
