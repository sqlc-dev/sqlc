package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type DropOwnedStmt struct {
	Roles    *ast.List
	Behavior DropBehavior
}

func (n *DropOwnedStmt) Pos() int {
	return 0
}
