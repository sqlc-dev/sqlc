package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type A_Expr struct {
	Kind     A_Expr_Kind
	Name     *ast.List
	Lexpr    ast.Node
	Rexpr    ast.Node
	Location int
}

func (n *A_Expr) Pos() int {
	return n.Location
}
