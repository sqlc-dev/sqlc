package ast

type AlterEventTrigStmt struct {
	Tag NodeTag[AlterEventTrigStmt] `json:"tag"`

	Trigname  *string
	Tgenabled byte
}

func (n *AlterEventTrigStmt) Pos() int {
	return 0
}
