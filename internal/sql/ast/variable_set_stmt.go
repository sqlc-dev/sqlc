package ast

type VariableSetStmt struct {
	Kind    VariableSetKind
	Name    *string
	Args    *List
	IsLocal bool
}

func (n *VariableSetStmt) Pos() int {
	return 0
}
