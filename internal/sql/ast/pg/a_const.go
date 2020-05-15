package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type A_Const struct {
	Val      ast.Node
	Location int
}

func (n *A_Const) Pos() int {
	return n.Location
}
