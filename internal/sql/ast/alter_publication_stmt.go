package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterPublicationStmt struct {
	Pubname      *string
	Options      *List
	Tables       *List
	ForAllTables bool
	TableAction  DefElemAction
}

func (n *AlterPublicationStmt) Pos() int {
	return 0
}
