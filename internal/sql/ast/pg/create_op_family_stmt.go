package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateOpFamilyStmt struct {
	Opfamilyname *ast.List
	Amname       *string
}

func (n *CreateOpFamilyStmt) Pos() int {
	return 0
}
