package ast

type GrantRoleStmt struct {
	Tag NodeTag[GrantRoleStmt] `json:"tag"`

	GrantedRoles *List `json:",omitempty"`
	GranteeRoles *List `json:",omitempty"`
	IsGrant      bool
	Grantor      *RoleSpec `json:",omitempty"`
	Behavior     DropBehavior
}

func (n *GrantRoleStmt) Pos() int {
	return 0
}
