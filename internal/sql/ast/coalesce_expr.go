package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type CoalesceExpr struct {
	Tag NodeTag[CoalesceExpr] `json:"tag"`

	Xpr            Node  `json:"xpr,omitempty"`
	Coalescetype   Oid   `json:"coalescetype"`
	Coalescecollid Oid   `json:"coalescecollid"`
	Args           *List `json:"args,omitempty"`
	Location       int   `json:"location"`
	// FirstArgType marks SQL Server's ISNULL, which takes its first
	// argument's type rather than the type its arguments combine to.
	FirstArgType bool `json:"first_arg_type,omitempty"`
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
