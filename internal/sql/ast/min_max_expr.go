package ast

type MinMaxExpr struct {
	Xpr          Node
	Minmaxtype   Oid
	Minmaxcollid Oid
	Inputcollid  Oid
	Op           MinMaxOp
	Args         *List
	Location     int
}

func (n *MinMaxExpr) Pos() int {
	return n.Location
}
