package ast

type A_Indices struct {
	IsSlice bool
	Lidx    Node
	Uidx    Node
}

func (n *A_Indices) Pos() int {
	return 0
}
