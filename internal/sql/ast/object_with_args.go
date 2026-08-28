package ast

type ObjectWithArgs struct {
	Tag NodeTag[ObjectWithArgs] `json:"tag"`

	Objname         *List `json:"objname,omitempty"`
	Objargs         *List `json:"objargs,omitempty"`
	ArgsUnspecified bool  `json:"args_unspecified"`
}

func (n *ObjectWithArgs) Pos() int {
	return 0
}
