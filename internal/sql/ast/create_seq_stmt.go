package ast

type CreateSeqStmt struct {
	Tag NodeTag[CreateSeqStmt] `json:"tag"`

	Sequence    *RangeVar `json:"sequence,omitempty"`
	Options     *List     `json:"options,omitempty"`
	OwnerId     Oid       `json:"owner_id"`
	ForIdentity bool      `json:"for_identity"`
	IfNotExists bool      `json:"if_not_exists"`
}

func (n *CreateSeqStmt) Pos() int {
	return 0
}
