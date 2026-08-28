package ast

type AlterTypeRenameValueStmt struct {
	Tag NodeTag[AlterTypeRenameValueStmt] `json:"tag"`

	Type     *TypeName
	OldValue *string
	NewValue *string
}

func (n *AlterTypeRenameValueStmt) Pos() int {
	return 0
}
