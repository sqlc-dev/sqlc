package ast

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

func (n *SubLink) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.astFormat(n.Testexpr)
	switch n.SubLinkType {
	case EXISTS_SUBLINK:
		buf.WriteString(" EXISTS (")
	case ANY_SUBLINK:
		buf.WriteString(" IN (")
	default:
		buf.WriteString(" (")
	}
	buf.astFormat(n.Subselect)
	buf.WriteString(")")
}
