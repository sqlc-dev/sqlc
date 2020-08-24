package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type ScalarArrayOpExpr struct {
	Xpr         Node
	Opno        Oid
	Opfuncid    Oid
	UseOr       bool
	Inputcollid Oid
	Args        *List
	Location    int
}

func (n *ScalarArrayOpExpr) Pos() int {
	return n.Location
}
