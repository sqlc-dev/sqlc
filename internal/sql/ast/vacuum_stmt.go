package ast

type VacuumStmt struct {
	Tag NodeTag[VacuumStmt] `json:"tag"`

	Options  int       `json:"options"`
	Relation *RangeVar `json:"relation,omitempty"`
	VaCols   *List     `json:"va_cols,omitempty"`
}

func (n *VacuumStmt) Pos() int {
	return 0
}
