package ast

type FuncExpr struct {
	Xpr            Node
	Funcid         Oid
	Funcresulttype Oid
	Funcretset     bool
	Funcvariadic   bool
	Funcformat     CoercionForm
	Funccollid     Oid
	Inputcollid    Oid
	Args           *List
	Location       int
}

func (n *FuncExpr) Pos() int {
	return n.Location
}
