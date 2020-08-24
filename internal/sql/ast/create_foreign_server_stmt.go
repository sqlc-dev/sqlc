package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateForeignServerStmt struct {
	Servername  *string
	Servertype  *string
	Version     *string
	Fdwname     *string
	IfNotExists bool
	Options     *ast.List
}

func (n *CreateForeignServerStmt) Pos() int {
	return 0
}
