package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type LockingClause struct {
	LockedRels *List
	Strength   LockClauseStrength
	WaitPolicy LockWaitPolicy
}

func (n *LockingClause) Pos() int {
	return 0
}
