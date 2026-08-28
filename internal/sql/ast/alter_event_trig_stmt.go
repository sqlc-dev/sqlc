package ast

type AlterEventTrigStmt struct {
	Tag NodeTag[AlterEventTrigStmt] `json:"tag"`

	Trigname  *string `json:",omitempty"`
	Tgenabled byte
}

func (n *AlterEventTrigStmt) Pos() int {
	return 0
}
