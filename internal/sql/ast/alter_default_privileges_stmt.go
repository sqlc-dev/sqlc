package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterDefaultPrivilegesStmt struct {
	Options *List
	Action  *GrantStmt
}

func (n *AlterDefaultPrivilegesStmt) Pos() int {
	return 0
}
