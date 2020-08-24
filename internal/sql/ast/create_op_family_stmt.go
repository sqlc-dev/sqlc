package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateOpFamilyStmt struct {
	Opfamilyname *List
	Amname       *string
}

func (n *CreateOpFamilyStmt) Pos() int {
	return 0
}
