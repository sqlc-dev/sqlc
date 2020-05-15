package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type TypeCast struct {
	Arg      ast.Node
	TypeName *TypeName
	Location int
}

func (n *TypeCast) Pos() int {
	return n.Location
}
