package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type MultiAssignRef struct {
	Source   ast.Node
	Colno    int
	Ncolumns int
}

func (n *MultiAssignRef) Pos() int {
	return 0
}
