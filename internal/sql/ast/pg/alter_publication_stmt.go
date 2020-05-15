package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterPublicationStmt struct {
	Pubname      *string
	Options      *ast.List
	Tables       *ast.List
	ForAllTables bool
	TableAction  DefElemAction
}

func (n *AlterPublicationStmt) Pos() int {
	return 0
}
