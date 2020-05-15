package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type RangeTableFuncCol struct {
	Colname       *string
	TypeName      *TypeName
	ForOrdinality bool
	IsNotNull     bool
	Colexpr       ast.Node
	Coldefexpr    ast.Node
	Location      int
}

func (n *RangeTableFuncCol) Pos() int {
	return n.Location
}
