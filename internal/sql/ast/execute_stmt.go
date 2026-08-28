package ast

type ExecuteStmt struct {
	Tag NodeTag[ExecuteStmt] `json:"tag"`

	Name   *string `json:"name,omitempty"`
	Params *List   `json:"params,omitempty"`
}

func (n *ExecuteStmt) Pos() int {
	return 0
}
