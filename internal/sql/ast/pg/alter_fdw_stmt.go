package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterFdwStmt struct {
	Fdwname     *string
	FuncOptions *ast.List
	Options     *ast.List
}

func (n *AlterFdwStmt) Pos() int {
	return 0
}
