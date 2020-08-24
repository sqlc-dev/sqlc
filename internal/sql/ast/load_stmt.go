package ast

type LoadStmt struct {
	Filename *string
}

func (n *LoadStmt) Pos() int {
	return 0
}
