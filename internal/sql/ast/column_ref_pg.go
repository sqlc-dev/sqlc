package ast

import ()

type ColumnRef struct {
	Fields   *List
	Location int
}

func (n *ColumnRef) Pos() int {
	return n.Location
}
