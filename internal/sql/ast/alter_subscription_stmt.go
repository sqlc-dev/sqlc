package ast

type AlterSubscriptionStmt struct {
	Tag NodeTag[AlterSubscriptionStmt] `json:"tag"`

	Kind        AlterSubscriptionType `json:"kind"`
	Subname     *string               `json:"subname,omitempty"`
	Conninfo    *string               `json:"conninfo,omitempty"`
	Publication *List                 `json:"publication,omitempty"`
	Options     *List                 `json:"options,omitempty"`
}

func (n *AlterSubscriptionStmt) Pos() int {
	return 0
}
