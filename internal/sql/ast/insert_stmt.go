package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type InsertStmt struct {
	Tag NodeTag[InsertStmt] `json:"tag"`

	Relation             *RangeVar             `json:"relation,omitempty"`
	Cols                 *List                 `json:"cols,omitempty"`
	SelectStmt           Node                  `json:"select_stmt,omitempty"`
	OnConflictClause     *OnConflictClause     `json:"on_conflict_clause,omitempty"`
	OnDuplicateKeyUpdate *OnDuplicateKeyUpdate `json:"on_duplicate_key_update,omitempty"` // MySQL-specific
	ReturningList        *List                 `json:"returning_list,omitempty"`
	WithClause           *WithClause           `json:"with_clause,omitempty"`
	Override             OverridingKind        `json:"override"`
	DefaultValues        bool                  `json:"default_values"` // SQLite-specific: INSERT INTO ... DEFAULT VALUES
	// PostgreSQL 18 RETURNING WITH (OLD AS ..., NEW AS ...) aliases
	ReturningOldAlias string `json:"returning_old_alias"`
	ReturningNewAlias string `json:"returning_new_alias"`
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
		// The boundary between the WITH clause and the statement's own
		// first keyword; the group keeps a break inside the WITH clause
		// from forcing the statement body apart clause by clause.
		buf.boundary(n.Relation)
		buf.Line()
		buf.Group()
		defer buf.EndGroup()
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
		if sel, ok := n.SelectStmt.(*SelectStmt); ok && items(sel.ValuesLists) {
			// The seam before a bare VALUES list keeps the author's own
			// choice — `) VALUES (` glued or VALUES on its own line — so it
			// gets a group of its own: a break inside either paren list must
			// not decide it. AttachComments reads the choice out of the
			// source (see valuesSeamBroken) and marks the boundary.
			buf.Group()
			buf.beforeClause(n.SelectStmt, d)
			buf.Line()
			buf.EndGroup()
		} else {
			buf.beforeClause(n.SelectStmt, d)
			buf.Line()
		}
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
