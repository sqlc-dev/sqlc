package ast

type CreateSubscriptionStmt struct {
	Tag NodeTag[CreateSubscriptionStmt] `json:"tag"`

	Subname     *string
	Conninfo    *string
	Publication *List
	Options     *List
}

func (n *CreateSubscriptionStmt) Pos() int {
	return 0
}
