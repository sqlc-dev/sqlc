package ast

type PrepareStmt struct {
	Name     *string
	Argtypes *List
	Query    Node
}

func (n *PrepareStmt) Pos() int {
	return 0
}
