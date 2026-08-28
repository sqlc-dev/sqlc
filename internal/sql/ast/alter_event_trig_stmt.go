package ast

type AlterEventTrigStmt struct {
	Tag NodeTag[AlterEventTrigStmt] `json:"tag"`

	Trigname  *string `json:"trigname,omitempty"`
	Tgenabled byte    `json:"tgenabled"`
}

func (n *AlterEventTrigStmt) Pos() int {
	return 0
}
