package ast

type FuncExpr struct {
	Tag NodeTag[FuncExpr] `json:"tag"`

	Xpr            Node         `json:"xpr,omitempty"`
	Funcid         Oid          `json:"funcid"`
	Funcresulttype Oid          `json:"funcresulttype"`
	Funcretset     bool         `json:"funcretset"`
	Funcvariadic   bool         `json:"funcvariadic"`
	Funcformat     CoercionForm `json:"funcformat"`
	Funccollid     Oid          `json:"funccollid"`
	Inputcollid    Oid          `json:"inputcollid"`
	Args           *List        `json:"args,omitempty"`
	Location       int          `json:"location"`
}

func (n *FuncExpr) Pos() int {
	return n.Location
}
