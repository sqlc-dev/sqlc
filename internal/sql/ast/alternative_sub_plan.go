package ast

type AlternativeSubPlan struct {
	Tag NodeTag[AlternativeSubPlan] `json:"tag"`

	Xpr      Node  `json:"xpr,omitempty"`
	Subplans *List `json:"subplans,omitempty"`
}

func (n *AlternativeSubPlan) Pos() int {
	return 0
}
