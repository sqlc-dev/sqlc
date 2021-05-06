package ast

type RenameTypeStmt struct {
	Type    *TypeName
	NewName *string
}

func (n *RenameTypeStmt) Pos() int {
	return 0
}
