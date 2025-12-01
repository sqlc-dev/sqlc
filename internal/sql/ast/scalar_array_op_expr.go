package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type ScalarArrayOpExpr struct {
	Xpr         Node
	Opno        Oid
	UseOr       bool
	Inputcollid Oid
	Args        *List
	Location    int
}

func (n *ScalarArrayOpExpr) Pos() int {
	return n.Location
}

func (n *ScalarArrayOpExpr) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	// ScalarArrayOpExpr represents "scalar op ANY/ALL (array)"
	// Args[0] is the left operand, Args[1] is the array
	if n.Args != nil && len(n.Args.Items) >= 2 {
		buf.astFormat(n.Args.Items[0], d)
		buf.WriteString(" = ") // TODO: Use actual operator based on Opno
		if n.UseOr {
			buf.WriteString("ANY(")
		} else {
			buf.WriteString("ALL(")
		}
		buf.astFormat(n.Args.Items[1], d)
		buf.WriteString(")")
	}
}
