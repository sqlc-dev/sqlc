package ast

type ReindexStmt struct {
	Tag NodeTag[ReindexStmt] `json:"tag"`

	Kind     ReindexObjectType
	Relation *RangeVar `json:",omitempty"`
	Name     *string   `json:",omitempty"`
	Options  int
}

func (n *ReindexStmt) Pos() int {
	return 0
}
