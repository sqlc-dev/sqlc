package ast

type RowMarkClause struct {
	Tag NodeTag[RowMarkClause] `json:"tag"`

	Rti        Index
	Strength   LockClauseStrength
	WaitPolicy LockWaitPolicy
	PushedDown bool
}

func (n *RowMarkClause) Pos() int {
	return 0
}
