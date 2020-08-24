package ast

import ()

type A_ArrayExpr struct {
	Elements *List
	Location int
}

func (n *A_ArrayExpr) Pos() int {
	return n.Location
}
