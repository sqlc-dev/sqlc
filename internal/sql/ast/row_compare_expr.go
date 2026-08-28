package ast

type RowCompareExpr struct {
	Tag NodeTag[RowCompareExpr] `json:"tag"`

	Xpr          Node `json:",omitempty"`
	Rctype       RowCompareType
	Opnos        *List `json:",omitempty"`
	Opfamilies   *List `json:",omitempty"`
	Inputcollids *List `json:",omitempty"`
	Largs        *List `json:",omitempty"`
	Rargs        *List `json:",omitempty"`
}

func (n *RowCompareExpr) Pos() int {
	return 0
}
