package ast

type UnlistenStmt struct {
	Tag NodeTag[UnlistenStmt] `json:"tag"`

	Conditionname *string
}

func (n *UnlistenStmt) Pos() int {
	return 0
}
