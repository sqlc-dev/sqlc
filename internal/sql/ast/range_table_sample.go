package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type RangeTableSample struct {
	Relation   Node
	Method     *List
	Args       *List
	Repeatable Node
	Location   int
}

func (n *RangeTableSample) Pos() int {
	return n.Location
}
