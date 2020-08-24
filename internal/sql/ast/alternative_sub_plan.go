package ast

import ()

type AlternativeSubPlan struct {
	Xpr      Node
	Subplans *List
}

func (n *AlternativeSubPlan) Pos() int {
	return 0
}
