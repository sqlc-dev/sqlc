package ast

import ()

type A_Const struct {
	Val      Node
	Location int
}

func (n *A_Const) Pos() int {
	return n.Location
}
