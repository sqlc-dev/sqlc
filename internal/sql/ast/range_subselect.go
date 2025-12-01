package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type RangeSubselect struct {
	Lateral  bool
	Subquery Node
	Alias    *Alias
}

func (n *RangeSubselect) Pos() int {
	return 0
}

func (n *RangeSubselect) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	if n.Lateral {
		buf.WriteString("LATERAL ")
	}
	buf.WriteString("(")
	buf.astFormat(n.Subquery, d)
	buf.WriteString(")")
	if n.Alias != nil {
		buf.WriteString(" AS ")
		buf.astFormat(n.Alias, d)
	}
}
