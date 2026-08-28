package ast

type ConvertRowtypeExpr struct {
	Tag NodeTag[ConvertRowtypeExpr] `json:"tag"`

	Xpr           Node
	Arg           Node
	Resulttype    Oid
	Convertformat CoercionForm
	Location      int
}

func (n *ConvertRowtypeExpr) Pos() int {
	return n.Location
}
