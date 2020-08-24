package ast

import ()

type TypeName struct {
	Names       *List
	TypeOid     Oid
	Setof       bool
	PctType     bool
	Typmods     *List
	Typemod     int32
	ArrayBounds *List
	Location    int
}

func (n *TypeName) Pos() int {
	return n.Location
}
