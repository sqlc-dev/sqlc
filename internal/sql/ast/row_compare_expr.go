package ast

type RowCompareExpr struct {
	Xpr          Node
	Rctype       RowCompareType
	Opnos        *List
	Opfamilies   *List
	Inputcollids *List
	Largs        *List
	Rargs        *List
}

func (n *RowCompareExpr) Pos() int {
	return 0
}
