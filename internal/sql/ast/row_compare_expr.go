package ast

type RowCompareExpr struct {
	Tag NodeTag[RowCompareExpr] `json:"tag"`

	Xpr          Node           `json:"xpr,omitempty"`
	Rctype       RowCompareType `json:"rctype"`
	Opnos        *List          `json:"opnos,omitempty"`
	Opfamilies   *List          `json:"opfamilies,omitempty"`
	Inputcollids *List          `json:"inputcollids,omitempty"`
	Largs        *List          `json:"largs,omitempty"`
	Rargs        *List          `json:"rargs,omitempty"`
}

func (n *RowCompareExpr) Pos() int {
	return 0
}
