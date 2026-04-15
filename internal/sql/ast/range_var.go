package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type RangeVar struct {
	Catalogname    *string
	Schemaname     *string
	Relname        *string
	Inh            bool
	Relpersistence byte
	Alias          *Alias
	Location       int
}

func (n *RangeVar) Pos() int {
	return n.Location
}

func (n *RangeVar) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	if n.Schemaname != nil && *n.Schemaname != "" {
		buf.WriteString(d.QuoteIdent(*n.Schemaname))
		buf.WriteString(".")
	}
	if n.Relname != nil {
		buf.WriteString(d.QuoteIdent(*n.Relname))
	}
	if n.Alias != nil {
		buf.WriteString(" AS ")
		buf.astFormat(n.Alias, d)
	}
}
