package ast

type A_Indirection struct {
	Tag NodeTag[A_Indirection] `json:"tag"`

	Arg         Node
	Indirection *List
}

func (n *A_Indirection) Pos() int {
	return 0
}
