package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type TruncateStmt struct {
	Relations   *ast.List
	RestartSeqs bool
	Behavior    DropBehavior
}

func (n *TruncateStmt) Pos() int {
	return 0
}
