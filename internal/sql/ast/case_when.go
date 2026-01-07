package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type CaseWhen struct {
	Xpr      Node
	Expr     Node
	Result   Node
	Location int
}

func (n *CaseWhen) Pos() int {
	return n.Location
}

func (n *CaseWhen) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("WHEN ")
	buf.astFormat(n.Expr, d)
	buf.WriteString(" THEN ")
	buf.astFormat(n.Result, d)
}
