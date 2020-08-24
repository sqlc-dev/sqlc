package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type GrantRoleStmt struct {
	GrantedRoles *List
	GranteeRoles *List
	IsGrant      bool
	AdminOpt     bool
	Grantor      *RoleSpec
	Behavior     DropBehavior
}

func (n *GrantRoleStmt) Pos() int {
	return 0
}
