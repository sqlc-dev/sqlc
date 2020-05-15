package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type PartitionElem struct {
	Name      *string
	Expr      ast.Node
	Collation *ast.List
	Opclass   *ast.List
	Location  int
}

func (n *PartitionElem) Pos() int {
	return n.Location
}
