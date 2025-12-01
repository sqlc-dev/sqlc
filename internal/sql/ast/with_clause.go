package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type WithClause struct {
	Ctes      *List
	Recursive bool
	Location  int
}

func (n *WithClause) Pos() int {
	return n.Location
}

func (n *WithClause) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("WITH ")
	if n.Recursive {
		buf.WriteString("RECURSIVE ")
	}
	buf.join(n.Ctes, d, ", ")
}
