package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type ViewStmt struct {
	View            *RangeVar
	Aliases         *List
	Query           ast.Node
	Replace         bool
	Options         *List
	WithCheckOption ViewCheckOption
}

func (n *ViewStmt) Pos() int {
	return 0
}
