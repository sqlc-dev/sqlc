package ast

type ArrayCoerceExpr struct {
	Tag NodeTag[ArrayCoerceExpr] `json:"tag"`

	Xpr          Node `json:",omitempty"`
	Arg          Node `json:",omitempty"`
	Elemfuncid   Oid
	Resulttype   Oid
	Resulttypmod int32
	Resultcollid Oid
	IsExplicit   bool
	Coerceformat CoercionForm
	Location     int
}

func (n *ArrayCoerceExpr) Pos() int {
	return n.Location
}
