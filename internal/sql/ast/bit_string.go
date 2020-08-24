package ast

type BitString struct {
	Str string
}

func (n *BitString) Pos() int {
	return 0
}
