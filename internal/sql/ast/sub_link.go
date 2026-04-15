package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type SubLinkType uint

const (
	EXISTS_SUBLINK SubLinkType = iota
	ALL_SUBLINK
	ANY_SUBLINK
	ROWCOMPARE_SUBLINK
	EXPR_SUBLINK
	MULTIEXPR_SUBLINK
	ARRAY_SUBLINK
	CTE_SUBLINK /* for SubPlans only */
)

type SubLink struct {
	Xpr         Node
	SubLinkType SubLinkType
	SubLinkId   int
	Testexpr    Node
	OperName    *List
	Subselect   Node
	Location    int
}

func (n *SubLink) Pos() int {
	return n.Location
}

func (n *SubLink) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	// Format the test expression if present (for IN subqueries etc.)
	hasTestExpr := n.Testexpr != nil
	if hasTestExpr {
		buf.astFormat(n.Testexpr, d)
	}
	switch n.SubLinkType {
	case EXISTS_SUBLINK:
		buf.WriteString("EXISTS (")
	case ANY_SUBLINK:
		if hasTestExpr {
			buf.WriteString(" IN (")
		} else {
			buf.WriteString("IN (")
		}
	default:
		if hasTestExpr {
			buf.WriteString(" (")
		} else {
			buf.WriteString("(")
		}
	}
	buf.astFormat(n.Subselect, d)
	buf.WriteString(")")
}
