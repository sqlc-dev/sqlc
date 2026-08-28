package ast

type AlterTableSetSchemaStmt struct {
	Tag NodeTag[AlterTableSetSchemaStmt] `json:"tag"`

	Table     *TableName
	NewSchema *string
	MissingOk bool
}

func (n *AlterTableSetSchemaStmt) Pos() int {
	return 0
}
