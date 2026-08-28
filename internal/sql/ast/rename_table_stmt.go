package ast

type RenameTableStmt struct {
	Tag NodeTag[RenameTableStmt] `json:"tag"`

	Table     *TableName `json:"table,omitempty"`
	NewName   *string    `json:"new_name,omitempty"`
	MissingOk bool       `json:"missing_ok"`
}

func (n *RenameTableStmt) Pos() int {
	return 0
}
