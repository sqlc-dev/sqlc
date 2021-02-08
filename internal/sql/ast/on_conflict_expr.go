package ast

type OnConflictExpr struct {
	Action          OnConflictAction
	ArbiterElems    *List
	ArbiterWhere    Node
	Constraint      Oid
	OnConflictSet   *List
	OnConflictWhere Node
	ExclRelIndex    int
	ExclRelTlist    *List
}

func (n *OnConflictExpr) Pos() int {
	return 0
}
