package ast

type AlternativeSubPlan struct {
	Tag NodeTag[AlternativeSubPlan] `json:"tag"`

	Xpr      Node  `json:",omitempty"`
	Subplans *List `json:",omitempty"`
}

func (n *AlternativeSubPlan) Pos() int {
	return 0
}
