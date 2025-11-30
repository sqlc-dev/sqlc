package ast

type BetweenExpr struct {
	// Expr is the value expression to be compared.
	Expr Node
	// Left is the left expression in the between statement.
	Left Node
	// Right is the right expression in the between statement.
	Right Node
	// Not is true, the expression is "not between".
	Not      bool
	Location int
}

func (n *BetweenExpr) Pos() int {
	return n.Location
}

func (n *BetweenExpr) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.astFormat(n.Expr)
	if n.Not {
		buf.WriteString(" NOT BETWEEN ")
	} else {
		buf.WriteString(" BETWEEN ")
	}
	buf.astFormat(n.Left)
	buf.WriteString(" AND ")
	buf.astFormat(n.Right)
}
