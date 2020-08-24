package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type ArrayExpr struct {
	Xpr           ast.Node
	ArrayTypeid   Oid
	ArrayCollid   Oid
	ElementTypeid Oid
	Elements      *ast.List
	Multidims     bool
	Location      int
}

func (n *ArrayExpr) Pos() int {
	return n.Location
}
