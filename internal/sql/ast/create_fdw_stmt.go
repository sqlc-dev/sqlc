package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateFdwStmt struct {
	Fdwname     *string
	FuncOptions *ast.List
	Options     *ast.List
}

func (n *CreateFdwStmt) Pos() int {
	return 0
}
