package ast

type ReindexStmt struct {
	Tag NodeTag[ReindexStmt] `json:"tag"`

	Kind     ReindexObjectType `json:"kind"`
	Relation *RangeVar         `json:"relation,omitempty"`
	Name     *string           `json:"name,omitempty"`
	Options  int               `json:"options"`
}

func (n *ReindexStmt) Pos() int {
	return 0
}
