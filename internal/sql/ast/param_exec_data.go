package ast

type ParamExecData struct {
	Tag NodeTag[ParamExecData] `json:"tag"`

	ExecPlan any   `json:"exec_plan,omitempty"`
	Value    Datum `json:"value"`
	Isnull   bool  `json:"isnull"`
}

func (n *ParamExecData) Pos() int {
	return 0
}
