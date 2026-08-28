package ast

type CreateSeqStmt struct {
	Tag NodeTag[CreateSeqStmt] `json:"tag"`

	Sequence    *RangeVar `json:",omitempty"`
	Options     *List     `json:",omitempty"`
	OwnerId     Oid
	ForIdentity bool
	IfNotExists bool
}

func (n *CreateSeqStmt) Pos() int {
	return 0
}
