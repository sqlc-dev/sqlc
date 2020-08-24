package ast

type AlterSystemStmt struct {
	Setstmt *VariableSetStmt
}

func (n *AlterSystemStmt) Pos() int {
	return 0
}
