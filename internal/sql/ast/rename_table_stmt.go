package ast

type RenameTableStmt struct {
	Table   *TableName
	NewName *string
}

func (n *RenameTableStmt) Pos() int {
	return 0
}
