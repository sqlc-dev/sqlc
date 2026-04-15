package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

// IntervalExpr represents a MySQL INTERVAL expression like "INTERVAL 1 DAY"
type IntervalExpr struct {
	Value    Node
	Unit     string
	Location int
}

func (n *IntervalExpr) Pos() int {
	return n.Location
}

func (n *IntervalExpr) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("INTERVAL ")
	buf.astFormat(n.Value, d)
	buf.WriteString(" ")
	buf.WriteString(n.Unit)
}
