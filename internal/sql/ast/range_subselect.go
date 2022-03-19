package ast

type RangeSubselect struct {
	Lateral  bool
	JoinType JoinType
	Subquery Node
	Alias    *Alias
}

func (n *RangeSubselect) Pos() int {
	return 0
}
