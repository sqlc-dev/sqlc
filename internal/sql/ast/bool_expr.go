package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type BoolExpr struct {
	Xpr      Node
	Boolop   BoolExprType
	Args     *List
	Location int
}

func (n *BoolExpr) Pos() int {
	return n.Location
}

func (n *BoolExpr) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	switch n.Boolop {
	case BoolExprTypeIsNull:
		if items(n.Args) && len(n.Args.Items) > 0 {
			buf.astFormat(n.Args.Items[0], d)
		}
		buf.WriteString(" IS NULL")
	case BoolExprTypeIsNotNull:
		if items(n.Args) && len(n.Args.Items) > 0 {
			buf.astFormat(n.Args.Items[0], d)
		}
		buf.WriteString(" IS NOT NULL")
	case BoolExprTypeNot:
		// NOT expression: format as NOT <arg>
		buf.WriteString("NOT ")
		if items(n.Args) && len(n.Args.Items) > 0 {
			buf.astFormat(n.Args.Items[0], d)
		}
	default:
		buf.WriteString("(")
		if items(n.Args) {
			switch n.Boolop {
			case BoolExprTypeAnd:
				buf.join(n.Args, d, " AND ")
			case BoolExprTypeOr:
				buf.join(n.Args, d, " OR ")
			}
		}
		buf.WriteString(")")
	}
}
