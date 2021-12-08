package ast

type CreateTableStmt struct {
	IfNotExists bool
	Name        *TableName
	Cols        []*ColumnDef
	ReferTable  *TableName
	Comment     string
	Inherits    []*TableName
}

func (n *CreateTableStmt) Pos() int {
	return 0
}
