package ast

type CompositeTypeStmt_PG struct {
	Typevar    *RangeVar
	Coldeflist *List
}

func (n *CompositeTypeStmt_PG) Pos() int {
	return 0
}
