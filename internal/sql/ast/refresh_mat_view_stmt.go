package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type RefreshMatViewStmt struct {
	Concurrent bool
	SkipData   bool
	Relation   *RangeVar
}

func (n *RefreshMatViewStmt) Pos() int {
	return 0
}

func (n *RefreshMatViewStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("REFRESH MATERIALIZED VIEW ")
	buf.astFormat(n.Relation, d)
}
