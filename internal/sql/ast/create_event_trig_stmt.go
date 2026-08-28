package ast

type CreateEventTrigStmt struct {
	Tag NodeTag[CreateEventTrigStmt] `json:"tag"`

	Trigname   *string
	Eventname  *string
	Whenclause *List
	Funcname   *List
}

func (n *CreateEventTrigStmt) Pos() int {
	return 0
}
