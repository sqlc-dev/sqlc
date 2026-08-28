package ast

type OnConflictExpr struct {
	Tag NodeTag[OnConflictExpr] `json:"tag"`

	Action          OnConflictAction
	ArbiterElems    *List `json:",omitempty"`
	ArbiterWhere    Node  `json:",omitempty"`
	Constraint      Oid
	OnConflictSet   *List `json:",omitempty"`
	OnConflictWhere Node  `json:",omitempty"`
	ExclRelIndex    int
	ExclRelTlist    *List `json:",omitempty"`
}

func (n *OnConflictExpr) Pos() int {
	return 0
}
