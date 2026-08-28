package ast

type CollateClause struct {
	Tag NodeTag[CollateClause] `json:"tag"`

	Arg      Node  `json:"arg,omitempty"`
	Collname *List `json:"collname,omitempty"`
	Location int   `json:"location"`
}

func (n *CollateClause) Pos() int {
	return n.Location
}
