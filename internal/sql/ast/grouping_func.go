package ast

type GroupingFunc struct {
	Xpr         Node
	Args        *List
	Refs        *List
	Cols        *List
	Agglevelsup Index
	Location    int
}

func (n *GroupingFunc) Pos() int {
	return n.Location
}
