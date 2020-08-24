package ast

import ()

type TruncateStmt struct {
	Relations   *List
	RestartSeqs bool
	Behavior    DropBehavior
}

func (n *TruncateStmt) Pos() int {
	return 0
}
