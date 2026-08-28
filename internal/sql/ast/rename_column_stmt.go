package ast

type RenameColumnStmt struct {
	Tag NodeTag[RenameColumnStmt] `json:"tag"`

	Table     *TableName `json:",omitempty"`
	Col       *ColumnRef `json:",omitempty"`
	NewName   *string    `json:",omitempty"`
	MissingOk bool
}

func (n *RenameColumnStmt) Pos() int {
	return 0
}
