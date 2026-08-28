package ast

type ParamExecData struct {
	Tag NodeTag[ParamExecData] `json:"tag"`

	ExecPlan any
	Value    Datum
	Isnull   bool
}

func (n *ParamExecData) Pos() int {
	return 0
}
