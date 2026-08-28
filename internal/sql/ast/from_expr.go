package ast

type FromExpr struct {
	Tag NodeTag[FromExpr] `json:"tag"`

	Fromlist *List `json:"fromlist,omitempty"`
	Quals    Node  `json:"quals,omitempty"`
}

func (n *FromExpr) Pos() int {
	return 0
}
