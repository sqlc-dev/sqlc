package ast

type TransactionStmt struct {
	Tag NodeTag[TransactionStmt] `json:"tag"`

	Kind    TransactionStmtKind
	Options *List
	Gid     *string
}

func (n *TransactionStmt) Pos() int {
	return 0
}
