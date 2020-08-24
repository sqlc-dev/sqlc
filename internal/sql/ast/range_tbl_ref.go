package ast

type RangeTblRef struct {
	Rtindex int
}

func (n *RangeTblRef) Pos() int {
	return 0
}
