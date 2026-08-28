package ast

type RowMarkClause struct {
	Tag NodeTag[RowMarkClause] `json:"tag"`

	Rti        Index              `json:"rti"`
	Strength   LockClauseStrength `json:"strength"`
	WaitPolicy LockWaitPolicy     `json:"wait_policy"`
	PushedDown bool               `json:"pushed_down"`
}

func (n *RowMarkClause) Pos() int {
	return 0
}
