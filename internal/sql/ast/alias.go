package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type Alias struct {
	Tag NodeTag[Alias] `json:"tag"`

	Aliasname *string `json:",omitempty"`
	Colnames  *List   `json:",omitempty"`
}

func (n *Alias) Pos() int {
	return 0
}

func (n *Alias) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	if n.Aliasname != nil {
		buf.WriteString(*n.Aliasname)
	}
	if items(n.Colnames) {
		buf.WriteString("(")
		buf.astFormat(n.Colnames, d)
		buf.WriteString(")")
	}
}
