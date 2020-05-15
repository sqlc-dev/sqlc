package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CaseWhen struct {
	Xpr      ast.Node
	Expr     ast.Node
	Result   ast.Node
	Location int
}

func (n *CaseWhen) Pos() int {
	return n.Location
}
