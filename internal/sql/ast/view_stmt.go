package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type ViewStmt struct {
	View            *RangeVar
	Aliases         *ast.List
	Query           ast.Node
	Replace         bool
	Options         *ast.List
	WithCheckOption ViewCheckOption
}

func (n *ViewStmt) Pos() int {
	return 0
}
