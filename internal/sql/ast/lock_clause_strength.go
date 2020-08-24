package pg

type LockClauseStrength uint

func (n *LockClauseStrength) Pos() int {
	return 0
}
