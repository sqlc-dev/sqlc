package pg

type LockWaitPolicy uint

func (n *LockWaitPolicy) Pos() int {
	return 0
}
