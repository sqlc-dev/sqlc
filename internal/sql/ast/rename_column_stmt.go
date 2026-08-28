package ast

type RenameColumnStmt struct {
	Tag NodeTag[RenameColumnStmt] `json:"tag"`

	Table     *TableName
	Col       *ColumnRef
	NewName   *string
	MissingOk bool
}

func (n *RenameColumnStmt) Pos() int {
	return 0
}
