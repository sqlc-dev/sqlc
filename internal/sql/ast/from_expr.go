package ast

type FromExpr struct {
	Tag NodeTag[FromExpr] `json:"tag"`

	Fromlist *List `json:",omitempty"`
	Quals    Node  `json:",omitempty"`
}

func (n *FromExpr) Pos() int {
	return 0
}
