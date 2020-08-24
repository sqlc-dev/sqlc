package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterOperatorStmt struct {
	Opername *ObjectWithArgs
	Options  *ast.List
}

func (n *AlterOperatorStmt) Pos() int {
	return 0
}
