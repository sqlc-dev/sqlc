package ast

type VacuumStmt struct {
	Options  int
	Relation *RangeVar
	VaCols   *List
}

func (n *VacuumStmt) Pos() int {
	return 0
}
