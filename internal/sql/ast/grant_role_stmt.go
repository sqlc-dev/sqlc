package ast

type GrantRoleStmt struct {
	Tag NodeTag[GrantRoleStmt] `json:"tag"`

	GrantedRoles *List        `json:"granted_roles,omitempty"`
	GranteeRoles *List        `json:"grantee_roles,omitempty"`
	IsGrant      bool         `json:"is_grant"`
	Grantor      *RoleSpec    `json:"grantor,omitempty"`
	Behavior     DropBehavior `json:"behavior"`
}

func (n *GrantRoleStmt) Pos() int {
	return 0
}
