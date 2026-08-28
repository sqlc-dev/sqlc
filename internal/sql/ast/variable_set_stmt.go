package ast

type VariableSetStmt struct {
	Tag NodeTag[VariableSetStmt] `json:"tag"`

	Kind    VariableSetKind
	Name    *string
	Args    *List
	IsLocal bool
}

func (n *VariableSetStmt) Pos() int {
	return 0
}
