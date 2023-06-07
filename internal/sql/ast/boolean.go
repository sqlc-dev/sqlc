package ast

type Boolean struct {
	Boolval bool
}

func (n *Boolean) Pos() int {
	return 0
}
