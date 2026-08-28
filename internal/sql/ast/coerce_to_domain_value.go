package ast

type CoerceToDomainValue struct {
	Tag NodeTag[CoerceToDomainValue] `json:"tag"`

	Xpr       Node `json:",omitempty"`
	TypeId    Oid
	TypeMod   int32
	Collation Oid
	Location  int
}

func (n *CoerceToDomainValue) Pos() int {
	return n.Location
}
