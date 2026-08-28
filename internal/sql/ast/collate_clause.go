package ast

type CollateClause struct {
	Tag NodeTag[CollateClause] `json:"tag"`

	Arg      Node
	Collname *List
	Location int
}

func (n *CollateClause) Pos() int {
	return n.Location
}
