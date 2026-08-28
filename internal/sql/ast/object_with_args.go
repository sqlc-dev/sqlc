package ast

type ObjectWithArgs struct {
	Tag NodeTag[ObjectWithArgs] `json:"tag"`

	Objname         *List `json:",omitempty"`
	Objargs         *List `json:",omitempty"`
	ArgsUnspecified bool
}

func (n *ObjectWithArgs) Pos() int {
	return 0
}
