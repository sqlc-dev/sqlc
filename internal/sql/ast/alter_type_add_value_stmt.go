package ast

type AlterTypeAddValueStmt struct {
	Type               *TypeName
	NewValue           *string
	SkipIfNewValExists bool
}

func (n *AlterTypeAddValueStmt) Pos() int {
	return 0
}
