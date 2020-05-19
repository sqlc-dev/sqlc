package ast

type DropTypeStmt struct {
	IfExists bool
	Types    []*TypeName
}

func (n *DropTypeStmt) Pos() int {
	return 0
}
