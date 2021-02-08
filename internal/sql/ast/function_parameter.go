package ast

type FunctionParameter struct {
	Name    *string
	ArgType *TypeName
	Mode    FunctionParameterMode
	Defexpr Node
}

func (n *FunctionParameter) Pos() int {
	return 0
}
