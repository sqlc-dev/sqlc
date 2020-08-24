package ast

import ()

type CoalesceExpr struct {
	Xpr            Node
	Coalescetype   Oid
	Coalescecollid Oid
	Args           *List
	Location       int
}

func (n *CoalesceExpr) Pos() int {
	return n.Location
}
