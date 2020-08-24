package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateDomainStmt struct {
	Domainname  *List
	TypeName    *TypeName
	CollClause  *CollateClause
	Constraints *List
}

func (n *CreateDomainStmt) Pos() int {
	return 0
}
