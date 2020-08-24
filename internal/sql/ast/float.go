package ast

type Float struct {
	Str string
}

func (n *Float) Pos() int {
	return 0
}
