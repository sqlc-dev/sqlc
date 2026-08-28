package ast

type CreateTrigStmt struct {
	Tag NodeTag[CreateTrigStmt] `json:"tag"`

	Trigname       *string   `json:"trigname,omitempty"`
	Relation       *RangeVar `json:"relation,omitempty"`
	Funcname       *List     `json:"funcname,omitempty"`
	Args           *List     `json:"args,omitempty"`
	Row            bool      `json:"row"`
	Timing         int16     `json:"timing"`
	Events         int16     `json:"events"`
	Columns        *List     `json:"columns,omitempty"`
	WhenClause     Node      `json:"when_clause,omitempty"`
	Isconstraint   bool      `json:"isconstraint"`
	TransitionRels *List     `json:"transition_rels,omitempty"`
	Deferrable     bool      `json:"deferrable"`
	Initdeferred   bool      `json:"initdeferred"`
	Constrrel      *RangeVar `json:"constrrel,omitempty"`
}

func (n *CreateTrigStmt) Pos() int {
	return 0
}
