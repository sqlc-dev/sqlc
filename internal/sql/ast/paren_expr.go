package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

// ParenExpr represents a parenthesized expression
type ParenExpr struct {
	Expr     Node
	Location int
}

func (n *ParenExpr) Pos() int {
	return n.Location
}

func (n *ParenExpr) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("(")
	buf.astFormat(n.Expr, d)
	buf.WriteString(")")
}
