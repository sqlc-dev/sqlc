package ast

import ()

type SortBy struct {
	Node        Node
	SortbyDir   SortByDir
	SortbyNulls SortByNulls
	UseOp       *List
	Location    int
}

func (n *SortBy) Pos() int {
	return n.Location
}
