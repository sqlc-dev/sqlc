package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type DropRoleStmt struct {
	Roles     *ast.List
	MissingOk bool
}

func (n *DropRoleStmt) Pos() int {
	return 0
}
