package ast

type CreateEventTrigStmt struct {
	Tag NodeTag[CreateEventTrigStmt] `json:"tag"`

	Trigname   *string `json:"trigname,omitempty"`
	Eventname  *string `json:"eventname,omitempty"`
	Whenclause *List   `json:"whenclause,omitempty"`
	Funcname   *List   `json:"funcname,omitempty"`
}

func (n *CreateEventTrigStmt) Pos() int {
	return 0
}
