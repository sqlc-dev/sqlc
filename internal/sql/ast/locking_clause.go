package ast

import ()

type LockingClause struct {
	LockedRels *List
	Strength   LockClauseStrength
	WaitPolicy LockWaitPolicy
}

func (n *LockingClause) Pos() int {
	return 0
}
