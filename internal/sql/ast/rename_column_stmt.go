package ast

type RenameColumnStmt struct {
	Table     *TableName
	Col       *ColumnRef
	NewName   *string
	MissingOk bool
}

func (n *RenameColumnStmt) Pos() int {
	return 0
}
