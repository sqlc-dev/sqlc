package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type A_Indices struct {
	IsSlice bool
	Lidx    Node
	Uidx    Node
}

func (n *A_Indices) Pos() int {
	return 0
}

func (n *A_Indices) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("[")
	if n.IsSlice {
		if set(n.Lidx) {
			buf.astFormat(n.Lidx, d)
		}
		buf.WriteString(":")
		if set(n.Uidx) {
			buf.astFormat(n.Uidx, d)
		}
	} else {
		buf.astFormat(n.Uidx, d)
	}
	buf.WriteString("]")
}
