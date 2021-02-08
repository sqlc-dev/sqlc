package ast

type CreateEventTrigStmt struct {
	Trigname   *string
	Eventname  *string
	Whenclause *List
	Funcname   *List
}

func (n *CreateEventTrigStmt) Pos() int {
	return 0
}
