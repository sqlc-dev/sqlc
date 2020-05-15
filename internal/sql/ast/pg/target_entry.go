package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type TargetEntry struct {
	Xpr             ast.Node
	Expr            ast.Node
	Resno           AttrNumber
	Resname         *string
	Ressortgroupref Index
	Resorigtbl      Oid
	Resorigcol      AttrNumber
	Resjunk         bool
}

func (n *TargetEntry) Pos() int {
	return 0
}
