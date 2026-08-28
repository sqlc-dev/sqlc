package ast

type A_Indirection struct {
	Tag NodeTag[A_Indirection] `json:"tag"`

	Arg         Node  `json:",omitempty"`
	Indirection *List `json:",omitempty"`
}

func (n *A_Indirection) Pos() int {
	return 0
}
