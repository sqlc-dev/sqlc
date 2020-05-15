package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
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
