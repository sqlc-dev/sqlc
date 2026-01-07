package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type InsertStmt struct {
	Relation             *RangeVar
	Cols                 *List
	SelectStmt           Node
	OnConflictClause     *OnConflictClause
	OnDuplicateKeyUpdate *OnDuplicateKeyUpdate // MySQL-specific
	ReturningList        *List
	WithClause           *WithClause
	Override             OverridingKind
	DefaultValues        bool // SQLite-specific: INSERT INTO ... DEFAULT VALUES
}

func (n *InsertStmt) Pos() int {
	return 0
}

func (n *InsertStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}

	if n.WithClause != nil {
		buf.astFormat(n.WithClause, d)
		buf.WriteString(" ")
	}

	buf.WriteString("INSERT INTO ")
	if n.Relation != nil {
		buf.astFormat(n.Relation, d)
	}
	if items(n.Cols) {
		buf.WriteString(" (")
		buf.astFormat(n.Cols, d)
		buf.WriteString(")")
	}

	if n.DefaultValues {
		buf.WriteString(" DEFAULT VALUES")
	} else if set(n.SelectStmt) {
		buf.WriteString(" ")
		buf.astFormat(n.SelectStmt, d)
	}

	if n.OnConflictClause != nil {
		buf.WriteString(" ")
		buf.astFormat(n.OnConflictClause, d)
	}

	if n.OnDuplicateKeyUpdate != nil {
		buf.WriteString(" ")
		buf.astFormat(n.OnDuplicateKeyUpdate, d)
	}

	if items(n.ReturningList) {
		buf.WriteString(" RETURNING ")
		buf.astFormat(n.ReturningList, d)
	}
}
