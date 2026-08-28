package ast

type CreateSubscriptionStmt struct {
	Tag NodeTag[CreateSubscriptionStmt] `json:"tag"`

	Subname     *string `json:"subname,omitempty"`
	Conninfo    *string `json:"conninfo,omitempty"`
	Publication *List   `json:"publication,omitempty"`
	Options     *List   `json:"options,omitempty"`
}

func (n *CreateSubscriptionStmt) Pos() int {
	return 0
}
