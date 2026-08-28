package ast

type TransactionStmt struct {
	Tag NodeTag[TransactionStmt] `json:"tag"`

	Kind    TransactionStmtKind `json:"kind"`
	Options *List               `json:"options,omitempty"`
	Gid     *string             `json:"gid,omitempty"`
}

func (n *TransactionStmt) Pos() int {
	return 0
}
