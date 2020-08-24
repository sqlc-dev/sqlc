package ast

type DropSubscriptionStmt struct {
	Subname   *string
	MissingOk bool
	Behavior  DropBehavior
}

func (n *DropSubscriptionStmt) Pos() int {
	return 0
}
