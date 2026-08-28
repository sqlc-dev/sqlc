package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type A_Const struct {
	Tag NodeTag[A_Const] `json:"tag"`

	Val      Node `json:",omitempty"`
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
