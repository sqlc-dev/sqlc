package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterUserMappingStmt struct {
	User       *RoleSpec
	Servername *string
	Options    *ast.List
}

func (n *AlterUserMappingStmt) Pos() int {
	return 0
}
