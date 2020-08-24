package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CollateClause struct {
	Arg      ast.Node
	Collname *ast.List
	Location int
}

func (n *CollateClause) Pos() int {
	return n.Location
}
