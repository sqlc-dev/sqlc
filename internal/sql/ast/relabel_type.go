package ast

type RelabelType struct {
	Tag NodeTag[RelabelType] `json:"tag"`

	Xpr           Node `json:",omitempty"`
	Arg           Node `json:",omitempty"`
	Resulttype    Oid
	Resulttypmod  int32
	Resultcollid  Oid
	Relabelformat CoercionForm
	Location      int
}

func (n *RelabelType) Pos() int {
	return n.Location
}
