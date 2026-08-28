package ast

type VariableShowStmt struct {
	Tag NodeTag[VariableShowStmt] `json:"tag"`

	Name *string `json:",omitempty"`
}

func (n *VariableShowStmt) Pos() int {
	return 0
}
