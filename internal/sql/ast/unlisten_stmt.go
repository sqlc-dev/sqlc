package ast

type UnlistenStmt struct {
	Tag NodeTag[UnlistenStmt] `json:"tag"`

	Conditionname *string `json:"conditionname,omitempty"`
}

func (n *UnlistenStmt) Pos() int {
	return 0
}
