package ast

import ()

type CompositeTypeStmt struct {
	Typevar    *RangeVar
	Coldeflist *List
}

func (n *CompositeTypeStmt) Pos() int {
	return 0
}
