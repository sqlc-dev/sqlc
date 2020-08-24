package ast

import ()

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
