package ast

type CreateEventTrigStmt struct {
	Tag NodeTag[CreateEventTrigStmt] `json:"tag"`

	Trigname   *string `json:",omitempty"`
	Eventname  *string `json:",omitempty"`
	Whenclause *List   `json:",omitempty"`
	Funcname   *List   `json:",omitempty"`
}

func (n *CreateEventTrigStmt) Pos() int {
	return 0
}
