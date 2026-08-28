package ast

type SetToDefault struct {
	Tag NodeTag[SetToDefault] `json:"tag"`

	Xpr       Node
	TypeId    Oid
	TypeMod   int32
	Collation Oid
	Location  int
}

func (n *SetToDefault) Pos() int {
	return n.Location
}
