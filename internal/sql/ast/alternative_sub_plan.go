package ast

type AlternativeSubPlan struct {
	Tag NodeTag[AlternativeSubPlan] `json:"tag"`

	Xpr      Node
	Subplans *List
}

func (n *AlternativeSubPlan) Pos() int {
	return 0
}
