package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type NamedArgExpr struct {
	Xpr       ast.Node
	Arg       ast.Node
	Name      *string
	Argnumber int
	Location  int
}

func (n *NamedArgExpr) Pos() int {
	return n.Location
}
