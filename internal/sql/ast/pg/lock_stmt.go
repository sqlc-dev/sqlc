package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type LockStmt struct {
	Relations *ast.List
	Mode      int
	Nowait    bool
}

func (n *LockStmt) Pos() int {
	return 0
}
