package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type RefreshMatViewStmt struct {
	Tag NodeTag[RefreshMatViewStmt] `json:"tag"`

	Concurrent bool      `json:"concurrent"`
	SkipData   bool      `json:"skip_data"`
	Relation   *RangeVar `json:"relation,omitempty"`
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
