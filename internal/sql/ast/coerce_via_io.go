package ast

type CoerceViaIO struct {
	Tag NodeTag[CoerceViaIO] `json:"tag"`

	Xpr          Node         `json:"xpr,omitempty"`
	Arg          Node         `json:"arg,omitempty"`
	Resulttype   Oid          `json:"resulttype"`
	Resultcollid Oid          `json:"resultcollid"`
	Coerceformat CoercionForm `json:"coerceformat"`
	Location     int          `json:"location"`
}

func (n *CoerceViaIO) Pos() int {
	return n.Location
}
