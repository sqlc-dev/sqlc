package ast

type RuleStmt struct {
	Tag NodeTag[RuleStmt] `json:"tag"`

	Relation    *RangeVar `json:",omitempty"`
	Rulename    *string   `json:",omitempty"`
	WhereClause Node      `json:",omitempty"`
	Event       CmdType
	Instead     bool
	Actions     *List `json:",omitempty"`
	Replace     bool
}

func (n *RuleStmt) Pos() int {
	return 0
}
