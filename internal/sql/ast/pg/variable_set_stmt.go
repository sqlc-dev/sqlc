package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type VariableSetStmt struct {
	Kind    VariableSetKind
	Name    *string
	Args    *ast.List
	IsLocal bool
}

func (n *VariableSetStmt) Pos() int {
	return 0
}
