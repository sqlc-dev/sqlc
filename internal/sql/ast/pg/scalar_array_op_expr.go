package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type ScalarArrayOpExpr struct {
	Xpr         ast.Node
	Opno        Oid
	Opfuncid    Oid
	UseOr       bool
	Inputcollid Oid
	Args        *ast.List
	Location    int
}

func (n *ScalarArrayOpExpr) Pos() int {
	return n.Location
}
