package ast

type ObjectWithArgs struct {
	Objname         *List
	Objargs         *List
	ArgsUnspecified bool
}

func (n *ObjectWithArgs) Pos() int {
	return 0
}
