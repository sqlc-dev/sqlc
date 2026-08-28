package ast

type MinMaxExpr struct {
	Tag NodeTag[MinMaxExpr] `json:"tag"`

	Xpr          Node `json:",omitempty"`
	Minmaxtype   Oid
	Minmaxcollid Oid
	Inputcollid  Oid
	Op           MinMaxOp
	Args         *List `json:",omitempty"`
	Location     int
}

func (n *MinMaxExpr) Pos() int {
	return n.Location
}
