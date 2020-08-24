package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type RangeTableSample struct {
	Relation   ast.Node
	Method     *List
	Args       *List
	Repeatable ast.Node
	Location   int
}

func (n *RangeTableSample) Pos() int {
	return n.Location
}
