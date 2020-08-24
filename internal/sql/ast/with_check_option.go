package ast

import ()

type WithCheckOption struct {
	Kind     WCOKind
	Relname  *string
	Polname  *string
	Qual     Node
	Cascaded bool
}

func (n *WithCheckOption) Pos() int {
	return 0
}
