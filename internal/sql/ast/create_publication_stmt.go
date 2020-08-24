package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreatePublicationStmt struct {
	Pubname      *string
	Options      *List
	Tables       *List
	ForAllTables bool
}

func (n *CreatePublicationStmt) Pos() int {
	return 0
}
