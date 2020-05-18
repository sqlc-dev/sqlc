package ast

type SelectStmt struct {
	Fields *List
	From   *List
}

func (n *SelectStmt) Pos() int {
	return 0
}
