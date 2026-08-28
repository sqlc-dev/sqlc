package ast

type DropSubscriptionStmt struct {
	Tag NodeTag[DropSubscriptionStmt] `json:"tag"`

	Subname   *string      `json:"subname,omitempty"`
	MissingOk bool         `json:"missing_ok"`
	Behavior  DropBehavior `json:"behavior"`
}

func (n *DropSubscriptionStmt) Pos() int {
	return 0
}
