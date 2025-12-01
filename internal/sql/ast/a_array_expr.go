package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type A_ArrayExpr struct {
	Elements *List
	Location int
}

func (n *A_ArrayExpr) Pos() int {
	return n.Location
}

func (n *A_ArrayExpr) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("ARRAY[")
	buf.join(n.Elements, d, ", ")
	buf.WriteString("]")
}
