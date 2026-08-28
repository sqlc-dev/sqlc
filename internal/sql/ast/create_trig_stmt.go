package ast

type CreateTrigStmt struct {
	Tag NodeTag[CreateTrigStmt] `json:"tag"`

	Trigname       *string   `json:",omitempty"`
	Relation       *RangeVar `json:",omitempty"`
	Funcname       *List     `json:",omitempty"`
	Args           *List     `json:",omitempty"`
	Row            bool
	Timing         int16
	Events         int16
	Columns        *List `json:",omitempty"`
	WhenClause     Node  `json:",omitempty"`
	Isconstraint   bool
	TransitionRels *List `json:",omitempty"`
	Deferrable     bool
	Initdeferred   bool
	Constrrel      *RangeVar `json:",omitempty"`
}

func (n *CreateTrigStmt) Pos() int {
	return 0
}
