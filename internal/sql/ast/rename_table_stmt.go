package ast

type RenameTableStmt struct {
	Table     *TableName
	NewName   *string
	MissingOk bool
}

func (n *RenameTableStmt) Pos() int {
	return 0
}
