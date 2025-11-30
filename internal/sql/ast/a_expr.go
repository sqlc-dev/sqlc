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
	switch n.Kind {
	case A_Expr_Kind_IN:
		buf.astFormat(n.Lexpr)
		buf.WriteString(" IN (")
		buf.astFormat(n.Rexpr)
		buf.WriteString(")")
	case A_Expr_Kind_LIKE:
		buf.astFormat(n.Lexpr)
		buf.WriteString(" LIKE ")
		buf.astFormat(n.Rexpr)
	case A_Expr_Kind_ILIKE:
		buf.astFormat(n.Lexpr)
		buf.WriteString(" ILIKE ")
		buf.astFormat(n.Rexpr)
	case A_Expr_Kind_SIMILAR:
		buf.astFormat(n.Lexpr)
		buf.WriteString(" SIMILAR TO ")
		buf.astFormat(n.Rexpr)
	case A_Expr_Kind_BETWEEN:
		buf.astFormat(n.Lexpr)
		buf.WriteString(" BETWEEN ")
		if l, ok := n.Rexpr.(*List); ok && len(l.Items) == 2 {
			buf.astFormat(l.Items[0])
			buf.WriteString(" AND ")
			buf.astFormat(l.Items[1])
		}
	case A_Expr_Kind_NOT_BETWEEN:
		buf.astFormat(n.Lexpr)
		buf.WriteString(" NOT BETWEEN ")
		if l, ok := n.Rexpr.(*List); ok && len(l.Items) == 2 {
			buf.astFormat(l.Items[0])
			buf.WriteString(" AND ")
			buf.astFormat(l.Items[1])
		}
	case A_Expr_Kind_DISTINCT:
		buf.astFormat(n.Lexpr)
		buf.WriteString(" IS DISTINCT FROM ")
		buf.astFormat(n.Rexpr)
	case A_Expr_Kind_NOT_DISTINCT:
		buf.astFormat(n.Lexpr)
		buf.WriteString(" IS NOT DISTINCT FROM ")
		buf.astFormat(n.Rexpr)
	case A_Expr_Kind_OP:
		// Standard binary operator
		if set(n.Lexpr) {
			buf.astFormat(n.Lexpr)
			buf.WriteString(" ")
		}
		buf.astFormat(n.Name)
		if set(n.Rexpr) {
			buf.WriteString(" ")
			buf.astFormat(n.Rexpr)
		}
	default:
		// Fallback for other cases
		if set(n.Lexpr) {
			buf.astFormat(n.Lexpr)
			buf.WriteString(" ")
		}
		buf.astFormat(n.Name)
		if set(n.Rexpr) {
			buf.WriteString(" ")
			buf.astFormat(n.Rexpr)
		}
	}
}
