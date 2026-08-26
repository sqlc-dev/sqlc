package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type CreateTableStmt struct {
	IfNotExists bool
	Name        *TableName
	Cols        []*ColumnDef
	ReferTable  *TableName
	Comment     string
	Inherits    []*TableName
	// Virtual marks a table backed by a module (SQLite's CREATE VIRTUAL
	// TABLE). The statement's argument list cannot be reconstructed from
	// this node, so it has no faithful rendering.
	Virtual bool
}

func (n *CreateTableStmt) Pos() int {
	return 0
}

func (n *CreateTableStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	// A virtual table cannot be printed back: its module arguments were
	// parsed away. Render nothing, which no verification accepts, so the
	// formatter keeps the statement as written.
	if n.Virtual {
		return
	}
	buf.WriteString("CREATE TABLE ")
	buf.astFormat(n.Name, d)

	buf.WriteString(" (")
	buf.Group()
	buf.Indent()
	if len(n.Cols) > 0 {
		buf.boundary(n.Cols[0])
	}
	buf.Softline()
	for i, col := range n.Cols {
		if i > 0 {
			buf.WriteString(",")
			buf.boundary(col)
			buf.Line()
		}
		buf.astFormat(col, d)
	}
	buf.EndIndent()
	buf.Softline()
	buf.EndGroup()
	buf.WriteString(")")
}
