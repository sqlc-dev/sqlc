package ast

type ReassignOwnedStmt struct {
	Tag NodeTag[ReassignOwnedStmt] `json:"tag"`

	Roles   *List     `json:"roles,omitempty"`
	Newrole *RoleSpec `json:"newrole,omitempty"`
}

func (n *ReassignOwnedStmt) Pos() int {
	return 0
}
