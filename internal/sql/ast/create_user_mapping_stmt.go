package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateUserMappingStmt struct {
	User        *RoleSpec
	Servername  *string
	IfNotExists bool
	Options     *List
}

func (n *CreateUserMappingStmt) Pos() int {
	return 0
}
