package ast

type RuleStmt struct {
	Tag NodeTag[RuleStmt] `json:"tag"`

	Relation    *RangeVar
	Rulename    *string
	WhereClause Node
	Event       CmdType
	Instead     bool
	Actions     *List
	Replace     bool
}

func (n *RuleStmt) Pos() int {
	return 0
}
