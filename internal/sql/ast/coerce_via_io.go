package ast

type CoerceViaIO struct {
	Tag NodeTag[CoerceViaIO] `json:"tag"`

	Xpr          Node `json:",omitempty"`
	Arg          Node `json:",omitempty"`
	Resulttype   Oid
	Resultcollid Oid
	Coerceformat CoercionForm
	Location     int
}

func (n *CoerceViaIO) Pos() int {
	return n.Location
}
