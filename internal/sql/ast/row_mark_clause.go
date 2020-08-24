package ast

type RowMarkClause struct {
	Rti        Index
	Strength   LockClauseStrength
	WaitPolicy LockWaitPolicy
	PushedDown bool
}

func (n *RowMarkClause) Pos() int {
	return 0
}
