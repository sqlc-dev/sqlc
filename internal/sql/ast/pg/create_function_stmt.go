package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateFunctionStmt struct {
	Replace    bool
	Funcname   *ast.List
	Parameters *ast.List
	ReturnType *TypeName
	Options    *ast.List
	WithClause *ast.List
}

func (n *CreateFunctionStmt) Pos() int {
	return 0
}
