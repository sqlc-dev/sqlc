package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type ResTarget struct {
	Name        *string
	Indirection *ast.List
	Val         ast.Node
	Location    int
}

func (n *ResTarget) Pos() int {
	return n.Location
}
