package ast

type AlterTableSetSchemaStmt struct {
	Table     *TableName
	NewSchema *string
	MissingOk bool
}

func (n *AlterTableSetSchemaStmt) Pos() int {
	return 0
}
