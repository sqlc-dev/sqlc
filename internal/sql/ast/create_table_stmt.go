package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type CreateTableStmt struct {
	IfNotExists bool
	Name        *TableName
	Cols        []*ColumnDef
	ReferTable  *TableName
	Comment     string
	Inherits    []*TableName
	// Incomplete marks a statement whose source carried syntax this node
	// does not model — a virtual table's module arguments, column or table
	// constraints beyond plain NOT NULL and PRIMARY KEY, table options,
	// TEMP, or an AS SELECT body. No faithful rendering exists for it.
	Incomplete bool
}

func (n *CreateTableStmt) Pos() int {
	return 0
}

func (n *CreateTableStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	// An incomplete statement cannot be printed back: part of its source
	// was parsed away. Render nothing, which no verification accepts, so
	// the formatter keeps the statement as written.
	if n.Incomplete {
		return
	}
	buf.WriteString("CREATE TABLE ")
	if n.IfNotExists {
		buf.WriteString("IF NOT EXISTS ")
	}
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
