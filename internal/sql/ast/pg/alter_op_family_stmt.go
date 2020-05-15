package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterOpFamilyStmt struct {
	Opfamilyname *ast.List
	Amname       *string
	IsDrop       bool
	Items        *ast.List
}

func (n *AlterOpFamilyStmt) Pos() int {
	return 0
}
