package ast

type RuleStmt struct {
	Tag NodeTag[RuleStmt] `json:"tag"`

	Relation    *RangeVar `json:"relation,omitempty"`
	Rulename    *string   `json:"rulename,omitempty"`
	WhereClause Node      `json:"where_clause,omitempty"`
	Event       CmdType   `json:"event"`
	Instead     bool      `json:"instead"`
	Actions     *List     `json:"actions,omitempty"`
	Replace     bool      `json:"replace"`
}

func (n *RuleStmt) Pos() int {
	return 0
}
