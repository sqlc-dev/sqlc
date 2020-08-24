package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CoerceViaIO struct {
	Xpr          ast.Node
	Arg          ast.Node
	Resulttype   Oid
	Resultcollid Oid
	Coerceformat CoercionForm
	Location     int
}

func (n *CoerceViaIO) Pos() int {
	return n.Location
}
