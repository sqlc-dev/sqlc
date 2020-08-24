package ast

type TableLikeClause struct {
	Relation *RangeVar
	Options  uint32
}

func (n *TableLikeClause) Pos() int {
	return 0
}
