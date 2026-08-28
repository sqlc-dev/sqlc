package ast

type ReassignOwnedStmt struct {
	Tag NodeTag[ReassignOwnedStmt] `json:"tag"`

	Roles   *List     `json:",omitempty"`
	Newrole *RoleSpec `json:",omitempty"`
}

func (n *ReassignOwnedStmt) Pos() int {
	return 0
}
