package ast

import ()

type DropOwnedStmt struct {
	Roles    *List
	Behavior DropBehavior
}

func (n *DropOwnedStmt) Pos() int {
	return 0
}
