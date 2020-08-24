package ast

type UnlistenStmt struct {
	Conditionname *string
}

func (n *UnlistenStmt) Pos() int {
	return 0
}
