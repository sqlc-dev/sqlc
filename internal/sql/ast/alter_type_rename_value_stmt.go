package ast

type AlterTypeRenameValueStmt struct {
	Tag NodeTag[AlterTypeRenameValueStmt] `json:"tag"`

	Type     *TypeName `json:",omitempty"`
	OldValue *string   `json:",omitempty"`
	NewValue *string   `json:",omitempty"`
}

func (n *AlterTypeRenameValueStmt) Pos() int {
	return 0
}
