package ast

type WithCheckOption struct {
	Tag NodeTag[WithCheckOption] `json:"tag"`

	Kind     WCOKind `json:"kind"`
	Relname  *string `json:"relname,omitempty"`
	Polname  *string `json:"polname,omitempty"`
	Qual     Node    `json:"qual,omitempty"`
	Cascaded bool    `json:"cascaded"`
}

func (n *WithCheckOption) Pos() int {
	return 0
}
