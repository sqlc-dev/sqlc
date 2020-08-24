package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CoerceViaIO struct {
	Xpr          Node
	Arg          Node
	Resulttype   Oid
	Resultcollid Oid
	Coerceformat CoercionForm
	Location     int
}

func (n *CoerceViaIO) Pos() int {
	return n.Location
}
