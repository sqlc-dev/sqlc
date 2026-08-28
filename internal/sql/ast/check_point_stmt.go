package ast

type CheckPointStmt struct {
	Tag NodeTag[CheckPointStmt] `json:"tag"`
}

func (n *CheckPointStmt) Pos() int {
	return 0
}
