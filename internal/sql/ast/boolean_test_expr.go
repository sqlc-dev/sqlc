package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type BooleanTest struct {
	Xpr          ast.Node
	Arg          ast.Node
	Booltesttype BoolTestType
	Location     int
}

func (n *BooleanTest) Pos() int {
	return n.Location
}
