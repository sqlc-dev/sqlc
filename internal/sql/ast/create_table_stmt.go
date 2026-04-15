package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

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

func (n *CreateTableStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("CREATE TABLE ")
	buf.astFormat(n.Name, d)

	buf.WriteString("(")
	for i, col := range n.Cols {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.astFormat(col, d)
	}
	buf.WriteString(")")
}
