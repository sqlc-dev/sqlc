package ast

type AlterTypeRenameValueStmt struct {
	Type     *TypeName
	OldValue *string
	NewValue *string
}

func (n *AlterTypeRenameValueStmt) Pos() int {
	return 0
}
