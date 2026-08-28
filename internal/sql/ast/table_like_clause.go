package ast

type TableLikeClause struct {
	Tag NodeTag[TableLikeClause] `json:"tag"`

	Relation *RangeVar
	Options  uint32
}

func (n *TableLikeClause) Pos() int {
	return 0
}
