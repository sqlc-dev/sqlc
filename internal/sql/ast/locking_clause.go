package ast

type LockingClause struct {
	LockedRels *List
	Strength   LockClauseStrength
	WaitPolicy LockWaitPolicy
}

func (n *LockingClause) Pos() int {
	return 0
}
