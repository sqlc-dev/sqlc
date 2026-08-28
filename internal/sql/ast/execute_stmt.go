package ast

type ExecuteStmt struct {
	Tag NodeTag[ExecuteStmt] `json:"tag"`

	Name   *string
	Params *List
}

func (n *ExecuteStmt) Pos() int {
	return 0
}
