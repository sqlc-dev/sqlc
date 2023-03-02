package ast

type FuncParamMode int

const (
	FuncParamIn FuncParamMode = iota
	FuncParamOut
	FuncParamInOut
	FuncParamVariadic
	FuncParamTable
	FuncParamDefault
)

type FuncParam struct {
	Name    *string
	Type    *TypeName
	DefExpr Node // Will always be &ast.TODO
	Mode    FuncParamMode
}

func (n *FuncParam) Pos() int {
	return 0
}
