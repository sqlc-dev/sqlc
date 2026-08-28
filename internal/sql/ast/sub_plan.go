package ast

type SubPlan struct {
	Tag NodeTag[SubPlan] `json:"tag"`

	Xpr               Node        `json:"xpr,omitempty"`
	SubLinkType       SubLinkType `json:"sub_link_type"`
	Testexpr          Node        `json:"testexpr,omitempty"`
	ParamIds          *List       `json:"param_ids,omitempty"`
	PlanId            int         `json:"plan_id"`
	PlanName          *string     `json:"plan_name,omitempty"`
	FirstColType      Oid         `json:"first_col_type"`
	FirstColTypmod    int32       `json:"first_col_typmod"`
	FirstColCollation Oid         `json:"first_col_collation"`
	UseHashTable      bool        `json:"use_hash_table"`
	UnknownEqFalse    bool        `json:"unknown_eq_false"`
	ParallelSafe      bool        `json:"parallel_safe"`
	SetParam          *List       `json:"set_param,omitempty"`
	ParParam          *List       `json:"par_param,omitempty"`
	Args              *List       `json:"args,omitempty"`
	StartupCost       Cost        `json:"startup_cost"`
	PerCallCost       Cost        `json:"per_call_cost"`
}

func (n *SubPlan) Pos() int {
	return 0
}
