package ast

type CreateStatsStmt struct {
	Tag NodeTag[CreateStatsStmt] `json:"tag"`

	Defnames    *List `json:",omitempty"`
	StatTypes   *List `json:",omitempty"`
	Exprs       *List `json:",omitempty"`
	Relations   *List `json:",omitempty"`
	IfNotExists bool
}

func (n *CreateStatsStmt) Pos() int {
	return 0
}
