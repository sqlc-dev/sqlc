package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type CoalesceExpr struct {
	Tag NodeTag[CoalesceExpr] `json:"tag"`

	Xpr            Node  `json:"xpr,omitempty"`
	Coalescetype   Oid   `json:"coalescetype"`
	Coalescecollid Oid   `json:"coalescecollid"`
	Args           *List `json:"args,omitempty"`
	Location       int   `json:"location"`
	// Func preserves the source identity of a function lowered to these
	// semantics. Its dialect decides any function-specific type rules.
	Func *FuncName `json:"func,omitempty"`
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
	if n.Func != nil {
		buf.astFormat(n.Func, d)
	} else {
		buf.WriteString("coalesce")
	}
	buf.WriteString("(")
	buf.astFormat(n.Args, d)
	buf.WriteString(")")
}
