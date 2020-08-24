package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type PrepareStmt struct {
	Name     *string
	Argtypes *ast.List
	Query    ast.Node
}

func (n *PrepareStmt) Pos() int {
	return 0
}
