package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type ArrayRef struct {
	Xpr             ast.Node
	Refarraytype    Oid
	Refelemtype     Oid
	Reftypmod       int32
	Refcollid       Oid
	Refupperindexpr *ast.List
	Reflowerindexpr *ast.List
	Refexpr         ast.Node
	Refassgnexpr    ast.Node
}

func (n *ArrayRef) Pos() int {
	return 0
}
