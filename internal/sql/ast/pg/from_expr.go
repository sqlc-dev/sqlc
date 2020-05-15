package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type FromExpr struct {
	Fromlist *ast.List
	Quals    ast.Node
}

func (n *FromExpr) Pos() int {
	return 0
}
