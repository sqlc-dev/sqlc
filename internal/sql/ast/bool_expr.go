package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type BoolExpr struct {
	Tag NodeTag[BoolExpr] `json:"tag"`

	Xpr      Node         `json:"xpr,omitempty"`
	Boolop   BoolExprType `json:"boolop"`
	Args     *List        `json:"args,omitempty"`
	Location int          `json:"location"`
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
		var op string
		switch n.Boolop {
		case BoolExprTypeAnd:
			op = "AND "
		case BoolExprTypeOr:
			op = "OR "
		}
		buf.WriteString("(")
		buf.Group()
		buf.Indent()
		buf.Softline()
		if items(n.Args) && op != "" {
			first := true
			for _, item := range n.Args.Items {
				if _, ok := item.(*TODO); ok {
					continue
				}
				if !first {
					buf.Line()
					buf.WriteString(op)
				}
				first = false
				buf.astFormat(item, d)
			}
		}
		buf.EndIndent()
		buf.Softline()
		buf.EndGroup()
		buf.WriteString(")")
	}
}
