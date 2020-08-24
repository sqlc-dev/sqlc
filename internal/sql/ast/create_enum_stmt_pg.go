package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateEnumStmt struct {
	TypeName *ast.List
	Vals     *ast.List
}

func (n *CreateEnumStmt) Pos() int {
	return 0
}
