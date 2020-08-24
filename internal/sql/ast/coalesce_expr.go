package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CoalesceExpr struct {
	Xpr            ast.Node
	Coalescetype   Oid
	Coalescecollid Oid
	Args           *ast.List
	Location       int
}

func (n *CoalesceExpr) Pos() int {
	return n.Location
}
