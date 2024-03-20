package ast

type OpExpr struct {
	Xpr          Node
	Opno         Oid
	Opresulttype Oid
	Opretset     bool
	Opcollid     Oid
	Inputcollid  Oid
	Args         *List
	Location     int
}

func (n *OpExpr) Pos() int {
	return n.Location
}
