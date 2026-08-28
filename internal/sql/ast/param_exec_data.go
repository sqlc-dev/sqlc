package ast

type ParamExecData struct {
	Tag NodeTag[ParamExecData] `json:"tag"`

	ExecPlan any `json:",omitempty"`
	Value    Datum
	Isnull   bool
}

func (n *ParamExecData) Pos() int {
	return 0
}
