package ast

type TransactionStmt struct {
	Tag NodeTag[TransactionStmt] `json:"tag"`

	Kind    TransactionStmtKind
	Options *List   `json:",omitempty"`
	Gid     *string `json:",omitempty"`
}

func (n *TransactionStmt) Pos() int {
	return 0
}
