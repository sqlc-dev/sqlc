package ast

type AlterSubscriptionStmt struct {
	Tag NodeTag[AlterSubscriptionStmt] `json:"tag"`

	Kind        AlterSubscriptionType
	Subname     *string `json:",omitempty"`
	Conninfo    *string `json:",omitempty"`
	Publication *List   `json:",omitempty"`
	Options     *List   `json:",omitempty"`
}

func (n *AlterSubscriptionStmt) Pos() int {
	return 0
}
