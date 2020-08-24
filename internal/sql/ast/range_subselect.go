package ast

import ()

type RangeSubselect struct {
	Lateral  bool
	Subquery Node
	Alias    *Alias
}

func (n *RangeSubselect) Pos() int {
	return 0
}
