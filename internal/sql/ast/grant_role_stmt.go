package ast

type GrantRoleStmt struct {
	Tag NodeTag[GrantRoleStmt] `json:"tag"`

	GrantedRoles *List
	GranteeRoles *List
	IsGrant      bool
	Grantor      *RoleSpec
	Behavior     DropBehavior
}

func (n *GrantRoleStmt) Pos() int {
	return 0
}
