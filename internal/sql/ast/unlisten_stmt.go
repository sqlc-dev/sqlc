package ast

type UnlistenStmt struct {
	Tag NodeTag[UnlistenStmt] `json:"tag"`

	Conditionname *string `json:",omitempty"`
}

func (n *UnlistenStmt) Pos() int {
	return 0
}
