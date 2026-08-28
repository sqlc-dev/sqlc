package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

// IntervalExpr represents a MySQL INTERVAL expression like "INTERVAL 1 DAY"
type IntervalExpr struct {
	Tag NodeTag[IntervalExpr] `json:"tag"`

	Value    Node   `json:"value,omitempty"`
	Unit     string `json:"unit"`
	Location int    `json:"location"`
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
