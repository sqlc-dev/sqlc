package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

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
	Xpr         ast.Node
	SubLinkType SubLinkType
	SubLinkId   int
	Testexpr    ast.Node
	OperName    *ast.List
	Subselect   ast.Node
	Location    int
}

func (n *SubLink) Pos() int {
	return n.Location
}
