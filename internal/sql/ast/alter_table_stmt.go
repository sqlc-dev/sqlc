package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type AlterTableStmt struct {
	Tag NodeTag[AlterTableStmt] `json:"tag"`

	// TODO: Only TableName or Relation should be defined
	Relation  *RangeVar  `json:"relation,omitempty"`
	Table     *TableName `json:"table,omitempty"`
	Cmds      *List      `json:"cmds,omitempty"`
	MissingOk bool       `json:"missing_ok"`
	Relkind   ObjectType `json:"relkind"`
	// Incomplete marks a statement whose source carried syntax this node
	// does not model, such as column constraints beyond plain NOT NULL and
	// PRIMARY KEY on an added column. No faithful rendering exists for it.
	Incomplete bool `json:"incomplete"`
}

func (n *AlterTableStmt) Pos() int {
	return 0
}

func (n *AlterTableStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	// An incomplete statement cannot be printed back: part of its source
	// was parsed away. Render nothing, which no verification accepts, so
	// the formatter keeps the statement as written.
	if n.Incomplete {
		return
	}
	buf.WriteString("ALTER TABLE ")
	buf.astFormat(n.Relation, d)
	buf.astFormat(n.Table, d)
	buf.astFormat(n.Cmds, d)
}
