package ast

type CaseTestExpr struct {
	Tag NodeTag[CaseTestExpr] `json:"tag"`

	Xpr       Node
	TypeId    Oid
	TypeMod   int32
	Collation Oid
}

func (n *CaseTestExpr) Pos() int {
	return 0
}
