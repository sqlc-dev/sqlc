package ast

// YDB specific
type UseStmt struct {
	Xpr      Node
	Location int
}

func (n *UseStmt) Pos() int {
	return n.Location
}
