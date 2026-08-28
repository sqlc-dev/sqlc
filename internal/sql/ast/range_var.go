package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type RangeVar struct {
	Tag NodeTag[RangeVar] `json:"tag"`

	Catalogname    *string `json:"catalogname,omitempty"`
	Schemaname     *string `json:"schemaname,omitempty"`
	Relname        *string `json:"relname,omitempty"`
	Inh            bool    `json:"inh"`
	Relpersistence byte    `json:"relpersistence"`
	Alias          *Alias  `json:"alias,omitempty"`
	Location       int     `json:"location"`
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
