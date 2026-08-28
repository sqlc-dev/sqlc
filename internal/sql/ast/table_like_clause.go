package ast

type TableLikeClause struct {
	Tag NodeTag[TableLikeClause] `json:"tag"`

	Relation *RangeVar `json:"relation,omitempty"`
	Options  uint32    `json:"options"`
}

func (n *TableLikeClause) Pos() int {
	return 0
}
