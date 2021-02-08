package ast

type DoStmt struct {
	Args *List
}

func (n *DoStmt) Pos() int {
	return 0
}
