package ast

type ConvertRowtypeExpr struct {
	Xpr           Node
	Arg           Node
	Resulttype    Oid
	Convertformat CoercionForm
	Location      int
}

func (n *ConvertRowtypeExpr) Pos() int {
	return n.Location
}
