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
	// PostgreSQL 18 RETURNING WITH (OLD AS ..., NEW AS ...) aliases
	ReturningOldAlias string
	ReturningNewAlias string
}

func (n *InsertStmt) Pos() int {
	return 0
}

func (n *InsertStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}

	buf.Group()
	defer buf.EndGroup()

	if n.WithClause != nil {
		buf.astFormat(n.WithClause, d)
		buf.Line()
	}

	buf.WriteString("INSERT INTO ")
	if n.Relation != nil {
		buf.astFormat(n.Relation, d)
	}
	if items(n.Cols) {
		buf.WriteString(" (")
		buf.Group()
		buf.Indent()
		buf.Softline()
		buf.joinComma(n.Cols, d)
		buf.EndIndent()
		buf.Softline()
		buf.EndGroup()
		buf.WriteString(")")
	}

	if n.DefaultValues {
		buf.WriteString(" DEFAULT VALUES")
	} else if set(n.SelectStmt) {
		buf.beforeClause(n.SelectStmt, d)
		buf.Line()
		buf.astFormat(n.SelectStmt, d)
	}

	if n.OnConflictClause != nil {
		buf.Line()
		buf.astFormat(n.OnConflictClause, d)
	}

	if n.OnDuplicateKeyUpdate != nil {
		buf.Line()
		buf.astFormat(n.OnDuplicateKeyUpdate, d)
	}

	if items(n.ReturningList) {
		buf.beforeClause(n.ReturningList, d)
		buf.Line()
		buf.WriteString("RETURNING ")
		formatReturningOptions(buf, d, n.ReturningOldAlias, n.ReturningNewAlias)
		buf.astFormat(n.ReturningList, d)
	}
}
