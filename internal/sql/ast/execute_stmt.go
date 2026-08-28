package ast

type ExecuteStmt struct {
	Tag NodeTag[ExecuteStmt] `json:"tag"`

	Name   *string `json:",omitempty"`
	Params *List   `json:",omitempty"`
}

func (n *ExecuteStmt) Pos() int {
	return 0
}
