package ast

type CaseTestExpr struct {
	Tag NodeTag[CaseTestExpr] `json:"tag"`

	Xpr       Node  `json:"xpr,omitempty"`
	TypeId    Oid   `json:"type_id"`
	TypeMod   int32 `json:"type_mod"`
	Collation Oid   `json:"collation"`
}

func (n *CaseTestExpr) Pos() int {
	return 0
}
