package ast

type NextValueExpr struct {
	Tag NodeTag[NextValueExpr] `json:"tag"`

	Xpr    Node `json:"xpr,omitempty"`
	Seqid  Oid  `json:"seqid"`
	TypeId Oid  `json:"type_id"`
}

func (n *NextValueExpr) Pos() int {
	return 0
}
