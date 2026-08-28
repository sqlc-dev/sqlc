package ast

type VariableShowStmt struct {
	Tag NodeTag[VariableShowStmt] `json:"tag"`

	Name *string
}

func (n *VariableShowStmt) Pos() int {
	return 0
}
