package ast

// ParenExpr represents a parenthesized expression
type ParenExpr struct {
	Expr     Node
	Location int
}

func (n *ParenExpr) Pos() int {
	return n.Location
}

func (n *ParenExpr) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString("(")
	buf.astFormat(n.Expr)
	buf.WriteString(")")
}
