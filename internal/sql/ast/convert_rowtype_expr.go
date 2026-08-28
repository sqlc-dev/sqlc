package ast

type ConvertRowtypeExpr struct {
	Tag NodeTag[ConvertRowtypeExpr] `json:"tag"`

	Xpr           Node         `json:"xpr,omitempty"`
	Arg           Node         `json:"arg,omitempty"`
	Resulttype    Oid          `json:"resulttype"`
	Convertformat CoercionForm `json:"convertformat"`
	Location      int          `json:"location"`
}

func (n *ConvertRowtypeExpr) Pos() int {
	return n.Location
}
