package ast

type CreateSeqStmt struct {
	Sequence    *RangeVar
	Options     *List
	OwnerId     Oid
	ForIdentity bool
	IfNotExists bool
}

func (n *CreateSeqStmt) Pos() int {
	return 0
}
