package ast

type CreateStatsStmt struct {
	Tag NodeTag[CreateStatsStmt] `json:"tag"`

	Defnames    *List `json:"defnames,omitempty"`
	StatTypes   *List `json:"stat_types,omitempty"`
	Exprs       *List `json:"exprs,omitempty"`
	Relations   *List `json:"relations,omitempty"`
	IfNotExists bool  `json:"if_not_exists"`
}

func (n *CreateStatsStmt) Pos() int {
	return 0
}
