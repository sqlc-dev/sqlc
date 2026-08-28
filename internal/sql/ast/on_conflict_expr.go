package ast

type OnConflictExpr struct {
	Tag NodeTag[OnConflictExpr] `json:"tag"`

	Action          OnConflictAction `json:"action"`
	ArbiterElems    *List            `json:"arbiter_elems,omitempty"`
	ArbiterWhere    Node             `json:"arbiter_where,omitempty"`
	Constraint      Oid              `json:"constraint"`
	OnConflictSet   *List            `json:"on_conflict_set,omitempty"`
	OnConflictWhere Node             `json:"on_conflict_where,omitempty"`
	ExclRelIndex    int              `json:"excl_rel_index"`
	ExclRelTlist    *List            `json:"excl_rel_tlist,omitempty"`
}

func (n *OnConflictExpr) Pos() int {
	return 0
}
