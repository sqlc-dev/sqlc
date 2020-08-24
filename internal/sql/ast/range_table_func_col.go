package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type RangeTableFuncCol struct {
	Colname       *string
	TypeName      *TypeName
	ForOrdinality bool
	IsNotNull     bool
	Colexpr       Node
	Coldefexpr    Node
	Location      int
}

func (n *RangeTableFuncCol) Pos() int {
	return n.Location
}
