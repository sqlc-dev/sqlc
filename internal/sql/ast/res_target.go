package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type ResTarget struct {
	Name        *string
	Indirection *List
	Val         Node
	Location    int
}

func (n *ResTarget) Pos() int {
	return n.Location
}

func (n *ResTarget) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	if set(n.Val) {
		buf.astFormat(n.Val, d)
		if n.Name != nil {
			buf.WriteString(" AS ")
			buf.WriteString(d.QuoteIdent(*n.Name))
		}
	} else {
		if n.Name != nil {
			buf.WriteString(d.QuoteIdent(*n.Name))
		}
	}
}
