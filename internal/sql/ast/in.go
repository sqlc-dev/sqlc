package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

// In describes a 'select foo in (bar, baz)' type statement, though there are multiple important variants handled.
type In struct {
	Tag NodeTag[In] `json:"tag"`

	// Expr is the value expression to be compared.
	Expr Node `json:"expr,omitempty"`
	// List is the list expression in compare list.
	List []Node `json:"list,omitempty"`
	// Not is true, the expression is "not in".
	Not bool `json:"not"`
	// Sel is the subquery, may be rewritten to other type of expression.
	Sel      Node `json:"sel,omitempty"`
	Location int  `json:"location"`
}

// Pos returns the location.
func (n *In) Pos() int {
	return n.Location
}

// Format formats the In expression.
func (n *In) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.astFormat(n.Expr, d)
	if n.Not {
		buf.WriteString(" NOT IN ")
	} else {
		buf.WriteString(" IN ")
	}
	if n.Sel != nil {
		buf.WriteString("(")
		buf.astFormat(n.Sel, d)
		buf.WriteString(")")
	} else if len(n.List) > 0 {
		buf.WriteString("(")
		for i, item := range n.List {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.astFormat(item, d)
		}
		buf.WriteString(")")
	}
}
