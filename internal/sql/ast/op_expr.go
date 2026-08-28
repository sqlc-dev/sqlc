package ast

type OpExpr struct {
	Tag NodeTag[OpExpr] `json:"tag"`

	Xpr          Node `json:",omitempty"`
	Opno         Oid
	Opresulttype Oid
	Opretset     bool
	Opcollid     Oid
	Inputcollid  Oid
	Args         *List `json:",omitempty"`
	Location     int
}

func (n *OpExpr) Pos() int {
	return n.Location
}
