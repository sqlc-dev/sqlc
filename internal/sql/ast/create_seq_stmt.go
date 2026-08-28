package ast

type CreateSeqStmt struct {
	Tag NodeTag[CreateSeqStmt] `json:"tag"`

	Sequence    *RangeVar
	Options     *List
	OwnerId     Oid
	ForIdentity bool
	IfNotExists bool
}

func (n *CreateSeqStmt) Pos() int {
	return 0
}
