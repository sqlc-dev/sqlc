package ast

type List struct {
	Items []Node
}

func (n *List) Pos() int {
	return 0
}
