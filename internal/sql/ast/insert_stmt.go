package ast

import "github.com/sqlc-dev/sqlc/internal/debug"

type InsertStmt struct {
	Relation         *RangeVar
	Cols             *List
	SelectStmt       Node
	OnConflictClause *OnConflictClause
	ReturningList    *List
	WithClause       *WithClause
	Override         OverridingKind
}

func (n *InsertStmt) Pos() int {
	return 0
}

func (n *InsertStmt) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}

	if n.WithClause != nil {
		buf.astFormat(n.WithClause)
		buf.WriteString(" ")
	}

	buf.WriteString("INSERT INTO ")
	if n.Relation != nil {
		debug.Dump(n.Relation)
		buf.astFormat(n.Relation)
	}
	if items(n.Cols) {
		buf.WriteString(" (")
		buf.astFormat(n.Cols)
		buf.WriteString(") ")
	}

	if set(n.SelectStmt) {
		buf.astFormat(n.SelectStmt)
	}
}
