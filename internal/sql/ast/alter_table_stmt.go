package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type AlterTableStmt struct {
	// TODO: Only TableName or Relation should be defined
	Relation  *RangeVar
	Table     *TableName
	Cmds      *List
	MissingOk bool
	Relkind   ObjectType
}

func (n *AlterTableStmt) Pos() int {
	return 0
}

func (n *AlterTableStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("ALTER TABLE ")
	buf.astFormat(n.Relation, d)
	buf.astFormat(n.Table, d)
	buf.astFormat(n.Cmds, d)
}
