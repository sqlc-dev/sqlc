package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type SortBy struct {
	Node        ast.Node
	SortbyDir   SortByDir
	SortbyNulls SortByNulls
	UseOp       *List
	Location    int
}

func (n *SortBy) Pos() int {
	return n.Location
}
