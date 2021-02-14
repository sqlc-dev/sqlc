package ast

type ExplainStmt struct {
	Query   Node
	Options *List
}

func (n *ExplainStmt) Pos() int {
	return 0
}
