package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateForeignTableStmt struct {
	Base       *CreateStmt
	Servername *string
	Options    *ast.List
}

func (n *CreateForeignTableStmt) Pos() int {
	return 0
}
