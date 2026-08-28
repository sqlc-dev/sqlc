package ast

type VacuumStmt struct {
	Tag NodeTag[VacuumStmt] `json:"tag"`

	Options  int
	Relation *RangeVar
	VaCols   *List
}

func (n *VacuumStmt) Pos() int {
	return 0
}
