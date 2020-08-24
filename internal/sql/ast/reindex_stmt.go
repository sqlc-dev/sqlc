package ast

type ReindexStmt struct {
	Kind     ReindexObjectType
	Relation *RangeVar
	Name     *string
	Options  int
}

func (n *ReindexStmt) Pos() int {
	return 0
}
