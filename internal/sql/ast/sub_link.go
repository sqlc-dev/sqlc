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
