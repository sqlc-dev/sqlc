package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateRoleStmt struct {
	StmtType RoleStmtType
	Role     *string
	Options  *ast.List
}

func (n *CreateRoleStmt) Pos() int {
	return 0
}
