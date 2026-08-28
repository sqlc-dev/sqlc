package ast

type ReassignOwnedStmt struct {
	Tag NodeTag[ReassignOwnedStmt] `json:"tag"`

	Roles   *List
	Newrole *RoleSpec
}

func (n *ReassignOwnedStmt) Pos() int {
	return 0
}
