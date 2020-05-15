package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateDomainStmt struct {
	Domainname  *ast.List
	TypeName    *TypeName
	CollClause  *CollateClause
	Constraints *ast.List
}

func (n *CreateDomainStmt) Pos() int {
	return 0
}
