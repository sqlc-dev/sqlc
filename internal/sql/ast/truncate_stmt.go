package ast

type TruncateStmt struct {
	Relations   *List
	RestartSeqs bool
	Behavior    DropBehavior
}

func (n *TruncateStmt) Pos() int {
	return 0
}
