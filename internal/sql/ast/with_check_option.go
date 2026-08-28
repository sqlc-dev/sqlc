package ast

type WithCheckOption struct {
	Tag NodeTag[WithCheckOption] `json:"tag"`

	Kind     WCOKind
	Relname  *string `json:",omitempty"`
	Polname  *string `json:",omitempty"`
	Qual     Node    `json:",omitempty"`
	Cascaded bool
}

func (n *WithCheckOption) Pos() int {
	return 0
}
