package ast

type NextValueExpr struct {
	Tag NodeTag[NextValueExpr] `json:"tag"`

	Xpr    Node `json:",omitempty"`
	Seqid  Oid
	TypeId Oid
}

func (n *NextValueExpr) Pos() int {
	return 0
}
