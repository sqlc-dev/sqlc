package ast

type AlterSubscriptionStmt struct {
	Kind        AlterSubscriptionType
	Subname     *string
	Conninfo    *string
	Publication *List
	Options     *List
}

func (n *AlterSubscriptionStmt) Pos() int {
	return 0
}
