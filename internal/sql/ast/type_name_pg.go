package ast

type TypeName_PG struct {
	Names       *List
	TypeOid     Oid
	Setof       bool
	PctType     bool
	Typmods     *List
	Typemod     int32
	ArrayBounds *List
	Location    int
}

func (n *TypeName_PG) Pos() int {
	return n.Location
}
