package ast

type DiscardStmt struct {
	Target DiscardMode
}

func (n *DiscardStmt) Pos() int {
	return 0
}
