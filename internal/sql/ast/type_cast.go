package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type TypeCast struct {
	Arg      Node
	TypeName *TypeName
	Location int
}

func (n *TypeCast) Pos() int {
	return n.Location
}
