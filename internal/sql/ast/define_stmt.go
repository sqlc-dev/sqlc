package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type DefineStmt struct {
	Kind        ObjectType
	Oldstyle    bool
	Defnames    *ast.List
	Args        *ast.List
	Definition  *ast.List
	IfNotExists bool
}

func (n *DefineStmt) Pos() int {
	return 0
}
