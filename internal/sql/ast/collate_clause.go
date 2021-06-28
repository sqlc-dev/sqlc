package ast

type CollateClause struct {
	Arg      Node
	Collname *List
	Location int
}

func (n *CollateClause) Pos() int {
	return n.Location
}
