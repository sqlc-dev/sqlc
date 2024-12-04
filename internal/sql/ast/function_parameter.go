package ast

type FunctionParameter struct {
	Name    *string
	ArgType *TypeName
	Mode    FunctionParameterMode
	Defexpr Node
	IsArray bool
}

func (n *FunctionParameter) Pos() int {
	return 0
}
