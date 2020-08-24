package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type GroupingFunc struct {
	Xpr         ast.Node
	Args        *ast.List
	Refs        *ast.List
	Cols        *ast.List
	Agglevelsup Index
	Location    int
}

func (n *GroupingFunc) Pos() int {
	return n.Location
}
