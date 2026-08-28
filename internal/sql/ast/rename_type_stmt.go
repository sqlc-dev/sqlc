package ast

type RenameTypeStmt struct {
	Tag NodeTag[RenameTypeStmt] `json:"tag"`

	Type    *TypeName `json:"type,omitempty"`
	NewName *string   `json:"new_name,omitempty"`
}

func (n *RenameTypeStmt) Pos() int {
	return 0
}
