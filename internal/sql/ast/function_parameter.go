package ast

type FunctionParameter struct {
	Tag NodeTag[FunctionParameter] `json:"tag"`

	Name    *string               `json:"name,omitempty"`
	ArgType *TypeName             `json:"arg_type,omitempty"`
	Mode    FunctionParameterMode `json:"mode"`
	Defexpr Node                  `json:"defexpr,omitempty"`
}

func (n *FunctionParameter) Pos() int {
	return 0
}
