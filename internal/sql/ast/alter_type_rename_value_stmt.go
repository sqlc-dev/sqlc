package ast

type AlterTypeRenameValueStmt struct {
	Tag NodeTag[AlterTypeRenameValueStmt] `json:"tag"`

	Type     *TypeName `json:"type,omitempty"`
	OldValue *string   `json:"old_value,omitempty"`
	NewValue *string   `json:"new_value,omitempty"`
}

func (n *AlterTypeRenameValueStmt) Pos() int {
	return 0
}
