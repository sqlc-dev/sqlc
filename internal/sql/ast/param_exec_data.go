package ast

type ParamExecData struct {
	ExecPlan any
	Value    Datum
	Isnull   bool
}

func (n *ParamExecData) Pos() int {
	return 0
}
