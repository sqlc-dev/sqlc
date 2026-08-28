package ast

type RenameTableStmt struct {
	Tag NodeTag[RenameTableStmt] `json:"tag"`

	Table     *TableName
	NewName   *string
	MissingOk bool
}

func (n *RenameTableStmt) Pos() int {
	return 0
}
