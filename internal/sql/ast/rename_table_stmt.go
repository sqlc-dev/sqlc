package ast

type RenameTableStmt struct {
	Tag NodeTag[RenameTableStmt] `json:"tag"`

	Table     *TableName `json:",omitempty"`
	NewName   *string    `json:",omitempty"`
	MissingOk bool
}

func (n *RenameTableStmt) Pos() int {
	return 0
}
