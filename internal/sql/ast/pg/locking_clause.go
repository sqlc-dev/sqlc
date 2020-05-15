package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type LockingClause struct {
	LockedRels *ast.List
	Strength   LockClauseStrength
	WaitPolicy LockWaitPolicy
}

func (n *LockingClause) Pos() int {
	return 0
}
