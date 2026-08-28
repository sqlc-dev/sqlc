package ast

type FunctionParameter struct {
	Tag NodeTag[FunctionParameter] `json:"tag"`

	Name    *string
	ArgType *TypeName
	Mode    FunctionParameterMode
	Defexpr Node
}

func (n *FunctionParameter) Pos() int {
	return 0
}
