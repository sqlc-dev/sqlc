package ast

type MinMaxExpr struct {
	Tag NodeTag[MinMaxExpr] `json:"tag"`

	Xpr          Node     `json:"xpr,omitempty"`
	Minmaxtype   Oid      `json:"minmaxtype"`
	Minmaxcollid Oid      `json:"minmaxcollid"`
	Inputcollid  Oid      `json:"inputcollid"`
	Op           MinMaxOp `json:"op"`
	Args         *List    `json:"args,omitempty"`
	Location     int      `json:"location"`
}

func (n *MinMaxExpr) Pos() int {
	return n.Location
}
