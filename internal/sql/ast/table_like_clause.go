package ast

type TableLikeClause struct {
	Tag NodeTag[TableLikeClause] `json:"tag"`

	Relation *RangeVar `json:",omitempty"`
	Options  uint32
}

func (n *TableLikeClause) Pos() int {
	return 0
}
