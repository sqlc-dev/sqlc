package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type NullTest struct {
	Xpr          ast.Node
	Arg          ast.Node
	Nulltesttype NullTestType
	Argisrow     bool
	Location     int
}

func (n *NullTest) Pos() int {
	return n.Location
}
