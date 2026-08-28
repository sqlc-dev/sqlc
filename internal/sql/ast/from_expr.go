package ast

type FromExpr struct {
	Tag NodeTag[FromExpr] `json:"tag"`

	Fromlist *List
	Quals    Node
}

func (n *FromExpr) Pos() int {
	return 0
}
