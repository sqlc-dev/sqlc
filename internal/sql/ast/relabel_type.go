package ast

type RelabelType struct {
	Tag NodeTag[RelabelType] `json:"tag"`

	Xpr           Node         `json:"xpr,omitempty"`
	Arg           Node         `json:"arg,omitempty"`
	Resulttype    Oid          `json:"resulttype"`
	Resulttypmod  int32        `json:"resulttypmod"`
	Resultcollid  Oid          `json:"resultcollid"`
	Relabelformat CoercionForm `json:"relabelformat"`
	Location      int          `json:"location"`
}

func (n *RelabelType) Pos() int {
	return n.Location
}
