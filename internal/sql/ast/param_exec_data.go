package ast

type ParamExecData struct {
	ExecPlan interface{}
	Value    Datum
	Isnull   bool
}

func (n *ParamExecData) Pos() int {
	return 0
}
