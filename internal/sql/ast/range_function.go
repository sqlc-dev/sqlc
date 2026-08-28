package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type RangeFunction struct {
	Tag NodeTag[RangeFunction] `json:"tag"`

	Lateral    bool   `json:"lateral"`
	Ordinality bool   `json:"ordinality"`
	IsRowsfrom bool   `json:"is_rowsfrom"`
	Functions  *List  `json:"functions,omitempty"`
	Alias      *Alias `json:"alias,omitempty"`
	Coldeflist *List  `json:"coldeflist,omitempty"`
}

func (n *RangeFunction) Pos() int {
	return 0
}

func (n *RangeFunction) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	if n.Lateral {
		buf.WriteString("LATERAL ")
	}
	buf.astFormat(n.Functions, d)
	if n.Ordinality {
		buf.WriteString(" WITH ORDINALITY")
	}
	if n.Alias != nil {
		buf.WriteString(" AS ")
		buf.astFormat(n.Alias, d)
	}
}
