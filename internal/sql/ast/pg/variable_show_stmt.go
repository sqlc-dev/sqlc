package pg

type VariableShowStmt struct {
	Name *string
}

func (n *VariableShowStmt) Pos() int {
	return 0
}
