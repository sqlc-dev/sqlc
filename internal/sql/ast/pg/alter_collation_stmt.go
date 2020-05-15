package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterCollationStmt struct {
	Collname *ast.List
}

func (n *AlterCollationStmt) Pos() int {
	return 0
}
