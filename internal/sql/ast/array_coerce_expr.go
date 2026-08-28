package ast

type ArrayCoerceExpr struct {
	Tag NodeTag[ArrayCoerceExpr] `json:"tag"`

	Xpr          Node
	Arg          Node
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
