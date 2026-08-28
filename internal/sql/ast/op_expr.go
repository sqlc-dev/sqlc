package ast

type OpExpr struct {
	Tag NodeTag[OpExpr] `json:"tag"`

	Xpr          Node  `json:"xpr,omitempty"`
	Opno         Oid   `json:"opno"`
	Opresulttype Oid   `json:"opresulttype"`
	Opretset     bool  `json:"opretset"`
	Opcollid     Oid   `json:"opcollid"`
	Inputcollid  Oid   `json:"inputcollid"`
	Args         *List `json:"args,omitempty"`
	Location     int   `json:"location"`
}

func (n *OpExpr) Pos() int {
	return n.Location
}
