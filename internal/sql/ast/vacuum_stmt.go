package ast

type VacuumStmt struct {
	Tag NodeTag[VacuumStmt] `json:"tag"`

	Options  int
	Relation *RangeVar `json:",omitempty"`
	VaCols   *List     `json:",omitempty"`
}

func (n *VacuumStmt) Pos() int {
	return 0
}
