package ast

type CreateStatsStmt struct {
	Tag NodeTag[CreateStatsStmt] `json:"tag"`

	Defnames    *List
	StatTypes   *List
	Exprs       *List
	Relations   *List
	IfNotExists bool
}

func (n *CreateStatsStmt) Pos() int {
	return 0
}
