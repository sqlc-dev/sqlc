package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterUserMappingStmt struct {
	User       *RoleSpec
	Servername *string
	Options    *List
}

func (n *AlterUserMappingStmt) Pos() int {
	return 0
}
