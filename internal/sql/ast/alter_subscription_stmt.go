package ast

type AlterSubscriptionStmt struct {
	Tag NodeTag[AlterSubscriptionStmt] `json:"tag"`

	Kind        AlterSubscriptionType
	Subname     *string
	Conninfo    *string
	Publication *List
	Options     *List
}

func (n *AlterSubscriptionStmt) Pos() int {
	return 0
}
