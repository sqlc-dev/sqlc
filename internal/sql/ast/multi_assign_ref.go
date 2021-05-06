package ast

type MultiAssignRef struct {
	Source   Node
	Colno    int
	Ncolumns int
}

func (n *MultiAssignRef) Pos() int {
	return 0
}
