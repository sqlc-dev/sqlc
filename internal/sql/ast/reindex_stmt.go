package ast

type ReindexStmt struct {
	Tag NodeTag[ReindexStmt] `json:"tag"`

	Kind     ReindexObjectType
	Relation *RangeVar
	Name     *string
	Options  int
}

func (n *ReindexStmt) Pos() int {
	return 0
}
