package ast

type MinMaxExpr struct {
	Tag NodeTag[MinMaxExpr] `json:"tag"`

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
