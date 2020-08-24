package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreatePublicationStmt struct {
	Pubname      *string
	Options      *ast.List
	Tables       *ast.List
	ForAllTables bool
}

func (n *CreatePublicationStmt) Pos() int {
	return 0
}
