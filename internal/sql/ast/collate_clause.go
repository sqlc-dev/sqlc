package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CollateClause struct {
	Arg      ast.Node
	Collname *List
	Location int
}

func (n *CollateClause) Pos() int {
	return n.Location
}
