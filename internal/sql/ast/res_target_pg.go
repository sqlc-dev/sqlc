package ast

type ResTarget_PG struct {
	Name        *string
	Indirection *List
	Val         Node
	Location    int
}

func (n *ResTarget_PG) Pos() int {
	return n.Location
}
