package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterFunctionStmt struct {
	Func    *ObjectWithArgs
	Actions *ast.List
}

func (n *AlterFunctionStmt) Pos() int {
	return 0
}
