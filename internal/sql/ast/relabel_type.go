package ast

type RelabelType struct {
	Xpr           Node
	Arg           Node
	Resulttype    Oid
	Resulttypmod  int32
	Resultcollid  Oid
	Relabelformat CoercionForm
	Location      int
}

func (n *RelabelType) Pos() int {
	return n.Location
}
