package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type A_Const struct {
	Val      Node
	Location int
}

func (n *A_Const) Pos() int {
	return n.Location
}

func (n *A_Const) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	if _, ok := n.Val.(*String); ok {
		buf.WriteString("'")
		buf.astFormat(n.Val, d)
		buf.WriteString("'")
	} else {
		buf.astFormat(n.Val, d)
	}
}
