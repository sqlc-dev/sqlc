package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type TypeName struct {
	Names       *ast.List
	TypeOid     Oid
	Setof       bool
	PctType     bool
	Typmods     *ast.List
	Typemod     int32
	ArrayBounds *ast.List
	Location    int
}

func (n *TypeName) Pos() int {
	return n.Location
}
