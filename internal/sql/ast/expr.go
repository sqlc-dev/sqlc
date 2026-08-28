package ast

type Expr struct {
	Tag NodeTag[Expr] `json:"tag"`
}

func (n *Expr) Pos() int {
	return 0
}
