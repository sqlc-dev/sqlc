package ast

type RenameTypeStmt struct {
	Tag NodeTag[RenameTypeStmt] `json:"tag"`

	Type    *TypeName
	NewName *string
}

func (n *RenameTypeStmt) Pos() int {
	return 0
}
