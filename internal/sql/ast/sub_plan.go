package ast

type SubPlan struct {
	Tag NodeTag[SubPlan] `json:"tag"`

	Xpr               Node `json:",omitempty"`
	SubLinkType       SubLinkType
	Testexpr          Node  `json:",omitempty"`
	ParamIds          *List `json:",omitempty"`
	PlanId            int
	PlanName          *string `json:",omitempty"`
	FirstColType      Oid
	FirstColTypmod    int32
	FirstColCollation Oid
	UseHashTable      bool
	UnknownEqFalse    bool
	ParallelSafe      bool
	SetParam          *List `json:",omitempty"`
	ParParam          *List `json:",omitempty"`
	Args              *List `json:",omitempty"`
	StartupCost       Cost
	PerCallCost       Cost
}

func (n *SubPlan) Pos() int {
	return 0
}
