package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterFunctionStmt struct {
	Func    *ObjectWithArgs
	Actions *List
}

func (n *AlterFunctionStmt) Pos() int {
	return 0
}
