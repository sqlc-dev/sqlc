package ast

type AlterTableSetSchemaStmt struct {
	Table     *TableName
	NewSchema *string
}

func (n *AlterTableSetSchemaStmt) Pos() int {
	return 0
}
