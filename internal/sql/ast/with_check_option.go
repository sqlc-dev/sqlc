package ast

type WithCheckOption struct {
	Tag NodeTag[WithCheckOption] `json:"tag"`

	Kind     WCOKind
	Relname  *string
	Polname  *string
	Qual     Node
	Cascaded bool
}

func (n *WithCheckOption) Pos() int {
	return 0
}
