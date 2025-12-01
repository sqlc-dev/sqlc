package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type LockingClause struct {
	LockedRels *List
	Strength   LockClauseStrength
	WaitPolicy LockWaitPolicy
}

func (n *LockingClause) Pos() int {
	return 0
}

// LockClauseStrength values (matching pg_query_go)
const (
	LockClauseStrengthUndefined      LockClauseStrength = 0
	LockClauseStrengthNone           LockClauseStrength = 1
	LockClauseStrengthForKeyShare    LockClauseStrength = 2
	LockClauseStrengthForShare       LockClauseStrength = 3
	LockClauseStrengthForNoKeyUpdate LockClauseStrength = 4
	LockClauseStrengthForUpdate      LockClauseStrength = 5
)

// LockWaitPolicy values
const (
	LockWaitPolicyBlock LockWaitPolicy = 1
	LockWaitPolicySkip  LockWaitPolicy = 2
	LockWaitPolicyError LockWaitPolicy = 3
)

func (n *LockingClause) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("FOR ")
	switch n.Strength {
	case LockClauseStrengthForKeyShare:
		buf.WriteString("KEY SHARE")
	case LockClauseStrengthForShare:
		buf.WriteString("SHARE")
	case LockClauseStrengthForNoKeyUpdate:
		buf.WriteString("NO KEY UPDATE")
	case LockClauseStrengthForUpdate:
		buf.WriteString("UPDATE")
	}
	if items(n.LockedRels) {
		buf.WriteString(" OF ")
		buf.join(n.LockedRels, d, ", ")
	}
	switch n.WaitPolicy {
	case LockWaitPolicySkip:
		buf.WriteString(" SKIP LOCKED")
	case LockWaitPolicyError:
		buf.WriteString(" NOWAIT")
	}
}
