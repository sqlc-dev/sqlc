package ast

type FunctionParameter struct {
	Tag NodeTag[FunctionParameter] `json:"tag"`

	Name    *string   `json:",omitempty"`
	ArgType *TypeName `json:",omitempty"`
	Mode    FunctionParameterMode
	Defexpr Node `json:",omitempty"`
}

func (n *FunctionParameter) Pos() int {
	return 0
}
