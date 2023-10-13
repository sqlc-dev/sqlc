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

func (n *CreateTableStmt) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString("CREATE TABLE ")
	buf.astFormat(n.Name)

	buf.WriteString("(")
	for i, col := range n.Cols {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.astFormat(col)
	}
	buf.WriteString(")")
}
