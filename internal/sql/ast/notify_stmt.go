package pg

type NotifyStmt struct {
	Conditionname *string
	Payload       *string
}

func (n *NotifyStmt) Pos() int {
	return 0
}
