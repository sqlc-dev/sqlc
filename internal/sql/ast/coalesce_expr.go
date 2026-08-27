package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type CoalesceExpr struct {
	Xpr            Node
	Coalescetype   Oid
	Coalescecollid Oid
	Args           *List
	Location       int
}

func (n *CoalesceExpr) Pos() int {
	return n.Location
}

func (n *CoalesceExpr) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	// Lower case, like every other function name: the printer upper-cases
	// keywords, and function names are identifiers, which fold lower.
	buf.WriteString("coalesce(")
	buf.astFormat(n.Args, d)
	buf.WriteString(")")
}
