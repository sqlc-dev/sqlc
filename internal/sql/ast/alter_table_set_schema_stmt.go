package ast

type AlterTableSetSchemaStmt struct {
	Tag NodeTag[AlterTableSetSchemaStmt] `json:"tag"`

	Table     *TableName `json:"table,omitempty"`
	NewSchema *string    `json:"new_schema,omitempty"`
	MissingOk bool       `json:"missing_ok"`
}

func (n *AlterTableSetSchemaStmt) Pos() int {
	return 0
}
