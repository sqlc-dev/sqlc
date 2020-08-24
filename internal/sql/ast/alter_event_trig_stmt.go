package ast

type AlterEventTrigStmt struct {
	Trigname  *string
	Tgenabled byte
}

func (n *AlterEventTrigStmt) Pos() int {
	return 0
}
