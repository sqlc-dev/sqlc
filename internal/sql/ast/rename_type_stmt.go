package ast

type RenameTypeStmt struct {
	Tag NodeTag[RenameTypeStmt] `json:"tag"`

	Type    *TypeName `json:",omitempty"`
	NewName *string   `json:",omitempty"`
}

func (n *RenameTypeStmt) Pos() int {
	return 0
}
