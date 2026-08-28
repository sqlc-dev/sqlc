package ast

type CreateSubscriptionStmt struct {
	Tag NodeTag[CreateSubscriptionStmt] `json:"tag"`

	Subname     *string `json:",omitempty"`
	Conninfo    *string `json:",omitempty"`
	Publication *List   `json:",omitempty"`
	Options     *List   `json:",omitempty"`
}

func (n *CreateSubscriptionStmt) Pos() int {
	return 0
}
