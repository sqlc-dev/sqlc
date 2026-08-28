package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type ResTarget struct {
	Tag NodeTag[ResTarget] `json:"tag"`

	Name        *string `json:"name,omitempty"`
	Indirection *List   `json:"indirection,omitempty"`
	Val         Node    `json:"val,omitempty"`
	Location    int     `json:"location"`
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
