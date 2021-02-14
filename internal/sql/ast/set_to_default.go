package ast

type SetToDefault struct {
	Xpr       Node
	TypeId    Oid
	TypeMod   int32
	Collation Oid
	Location  int
}

func (n *SetToDefault) Pos() int {
	return n.Location
}
