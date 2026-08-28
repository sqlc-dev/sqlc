package ast

type ObjectWithArgs struct {
	Tag NodeTag[ObjectWithArgs] `json:"tag"`

	Objname         *List
	Objargs         *List
	ArgsUnspecified bool
}

func (n *ObjectWithArgs) Pos() int {
	return 0
}
