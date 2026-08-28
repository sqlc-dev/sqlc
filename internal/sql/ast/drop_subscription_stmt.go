package ast

type DropSubscriptionStmt struct {
	Tag NodeTag[DropSubscriptionStmt] `json:"tag"`

	Subname   *string `json:",omitempty"`
	MissingOk bool
	Behavior  DropBehavior
}

func (n *DropSubscriptionStmt) Pos() int {
	return 0
}
