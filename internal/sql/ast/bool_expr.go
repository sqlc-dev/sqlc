package ast

type BoolExpr struct {
	Xpr      Node
	Boolop   BoolExprType
	Args     *List
	Location int
}

func (n *BoolExpr) Pos() int {
	return n.Location
}

func (n *BoolExpr) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	switch n.Boolop {
	case BoolExprTypeIsNull:
		if items(n.Args) && len(n.Args.Items) > 0 {
			buf.astFormat(n.Args.Items[0])
		}
		buf.WriteString(" IS NULL")
	case BoolExprTypeIsNotNull:
		if items(n.Args) && len(n.Args.Items) > 0 {
			buf.astFormat(n.Args.Items[0])
		}
		buf.WriteString(" IS NOT NULL")
	default:
		buf.WriteString("(")
		if items(n.Args) {
			switch n.Boolop {
			case BoolExprTypeAnd:
				buf.join(n.Args, " AND ")
			case BoolExprTypeOr:
				buf.join(n.Args, " OR ")
			case BoolExprTypeNot:
				buf.WriteString(" NOT ")
				buf.astFormat(n.Args)
			}
		}
		buf.WriteString(")")
	}
}
