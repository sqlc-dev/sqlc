package ast

type RenameColumnStmt struct {
	Tag NodeTag[RenameColumnStmt] `json:"tag"`

	Table     *TableName `json:"table,omitempty"`
	Col       *ColumnRef `json:"col,omitempty"`
	NewName   *string    `json:"new_name,omitempty"`
	MissingOk bool       `json:"missing_ok"`
}

func (n *RenameColumnStmt) Pos() int {
	return 0
}
