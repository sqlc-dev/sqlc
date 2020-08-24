package ast

import ()

type FromExpr struct {
	Fromlist *List
	Quals    Node
}

func (n *FromExpr) Pos() int {
	return 0
}
