package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type ExecuteStmt struct {
	Name   *string
	Params *ast.List
}

func (n *ExecuteStmt) Pos() int {
	return 0
}
