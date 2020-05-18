package ast

type CreateTableStmt struct {
	IfNotExists bool
	Name        *TableName
	Cols        []*ColumnDef
}

func (n *CreateTableStmt) Pos() int {
	return 0
}
