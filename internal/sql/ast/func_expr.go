package ast

type FuncExpr struct {
	Tag NodeTag[FuncExpr] `json:"tag"`

	Xpr            Node `json:",omitempty"`
	Funcid         Oid
	Funcresulttype Oid
	Funcretset     bool
	Funcvariadic   bool
	Funcformat     CoercionForm
	Funccollid     Oid
	Inputcollid    Oid
	Args           *List `json:",omitempty"`
	Location       int
}

func (n *FuncExpr) Pos() int {
	return n.Location
}
