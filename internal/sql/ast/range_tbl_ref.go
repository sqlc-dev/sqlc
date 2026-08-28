package ast

type RangeTblRef struct {
	Tag NodeTag[RangeTblRef] `json:"tag"`

	Rtindex int
}

func (n *RangeTblRef) Pos() int {
	return 0
}
