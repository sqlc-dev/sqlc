package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type RangeVar struct {
	Tag NodeTag[RangeVar] `json:"tag"`

	Catalogname    *string `json:",omitempty"`
	Schemaname     *string `json:",omitempty"`
	Relname        *string `json:",omitempty"`
	Inh            bool
	Relpersistence byte
	Alias          *Alias `json:",omitempty"`
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
