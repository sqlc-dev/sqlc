package ast

type Statement struct {
	Raw *RawStmt
}

func (n *Statement) Pos() int {
	return 0
}
