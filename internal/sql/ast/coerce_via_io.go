package ast

type CoerceViaIO struct {
	Tag NodeTag[CoerceViaIO] `json:"tag"`

	Xpr          Node
	Arg          Node
	Resulttype   Oid
	Resultcollid Oid
	Coerceformat CoercionForm
	Location     int
}

func (n *CoerceViaIO) Pos() int {
	return n.Location
}
