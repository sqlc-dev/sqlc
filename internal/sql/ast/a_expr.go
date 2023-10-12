package ast

type A_Expr struct {
	Kind     A_Expr_Kind
	Name     *List
	Lexpr    Node
	Rexpr    Node
	Location int
}

func (n *A_Expr) Pos() int {
	return n.Location
}

func (n *A_Expr) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.astFormat(n.Lexpr)
	buf.WriteString(" ")
	switch n.Kind {
	case A_Expr_Kind_IN:
		buf.WriteString(" IN (")
		buf.astFormat(n.Rexpr)
		buf.WriteString(")")
	case A_Expr_Kind_LIKE:
		buf.WriteString(" LIKE ")
		buf.astFormat(n.Rexpr)
	default:
		buf.astFormat(n.Name)
		buf.WriteString(" ")
		buf.astFormat(n.Rexpr)
	}
}
