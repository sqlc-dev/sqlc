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
	Tag NodeTag[SubLink] `json:"tag"`

	Xpr         Node        `json:"xpr,omitempty"`
	SubLinkType SubLinkType `json:"sub_link_type"`
	SubLinkId   int         `json:"sub_link_id"`
	Testexpr    Node        `json:"testexpr,omitempty"`
	OperName    *List       `json:"oper_name,omitempty"`
	Subselect   Node        `json:"subselect,omitempty"`
	Location    int         `json:"location"`
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
	buf.Group()
	buf.Indent()
	buf.Softline()
	buf.astFormat(n.Subselect, d)
	buf.EndIndent()
	buf.Softline()
	buf.EndGroup()
	buf.WriteString(")")
}
