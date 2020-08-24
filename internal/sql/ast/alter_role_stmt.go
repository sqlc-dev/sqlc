package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterRoleStmt struct {
	Role    *RoleSpec
	Options *List
	Action  int
}

func (n *AlterRoleStmt) Pos() int {
	return 0
}
