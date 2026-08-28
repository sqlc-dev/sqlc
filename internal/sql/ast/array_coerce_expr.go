package ast

type ArrayCoerceExpr struct {
	Tag NodeTag[ArrayCoerceExpr] `json:"tag"`

	Xpr          Node         `json:"xpr,omitempty"`
	Arg          Node         `json:"arg,omitempty"`
	Elemfuncid   Oid          `json:"elemfuncid"`
	Resulttype   Oid          `json:"resulttype"`
	Resulttypmod int32        `json:"resulttypmod"`
	Resultcollid Oid          `json:"resultcollid"`
	IsExplicit   bool         `json:"is_explicit"`
	Coerceformat CoercionForm `json:"coerceformat"`
	Location     int          `json:"location"`
}

func (n *ArrayCoerceExpr) Pos() int {
	return n.Location
}
