package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type SortBy struct {
	Node        ast.Node
	SortbyDir   SortByDir
	SortbyNulls SortByNulls
	UseOp       *ast.List
	Location    int
}

func (n *SortBy) Pos() int {
	return n.Location
}
