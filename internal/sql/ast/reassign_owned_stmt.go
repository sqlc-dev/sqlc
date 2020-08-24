package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type ReassignOwnedStmt struct {
	Roles   *ast.List
	Newrole *RoleSpec
}

func (n *ReassignOwnedStmt) Pos() int {
	return 0
}
