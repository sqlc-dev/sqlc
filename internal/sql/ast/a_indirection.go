package ast

type A_Indirection struct {
	Arg         Node
	Indirection *List
}

func (n *A_Indirection) Pos() int {
	return 0
}
