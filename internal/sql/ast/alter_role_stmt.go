package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterRoleStmt struct {
	Role    *RoleSpec
	Options *ast.List
	Action  int
}

func (n *AlterRoleStmt) Pos() int {
	return 0
}
