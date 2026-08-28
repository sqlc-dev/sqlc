package ast

type AlterTableSetSchemaStmt struct {
	Tag NodeTag[AlterTableSetSchemaStmt] `json:"tag"`

	Table     *TableName `json:",omitempty"`
	NewSchema *string    `json:",omitempty"`
	MissingOk bool
}

func (n *AlterTableSetSchemaStmt) Pos() int {
	return 0
}
