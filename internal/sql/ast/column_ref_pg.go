package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type ColumnRef struct {
	Fields   *ast.List
	Location int
}

func (n *ColumnRef) Pos() int {
	return n.Location
}
