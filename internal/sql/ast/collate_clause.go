package ast

type CollateClause struct {
	Tag NodeTag[CollateClause] `json:"tag"`

	Arg      Node  `json:",omitempty"`
	Collname *List `json:",omitempty"`
	Location int
}

func (n *CollateClause) Pos() int {
	return n.Location
}
