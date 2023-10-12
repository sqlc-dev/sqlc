package ast

type LockingClause struct {
	LockedRels *List
	Strength   LockClauseStrength
	WaitPolicy LockWaitPolicy
}

func (n *LockingClause) Pos() int {
	return 0
}

func (n *LockingClause) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString("FOR ")
	switch n.Strength {
	case 3:
		buf.WriteString("SHARE")
	case 5:
		buf.WriteString("UPDATE")
	}
}
