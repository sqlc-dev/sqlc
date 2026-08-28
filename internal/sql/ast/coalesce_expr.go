package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type CoalesceExpr struct {
	Tag NodeTag[CoalesceExpr] `json:"tag"`

	Xpr            Node `json:",omitempty"`
	Coalescetype   Oid
	Coalescecollid Oid
	Args           *List `json:",omitempty"`
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
