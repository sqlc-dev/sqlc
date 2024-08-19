package ast

type GrantRoleStmt struct {
	GrantedRoles *List
	GranteeRoles *List
	IsGrant      bool
	Grantor      *RoleSpec
	Behavior     DropBehavior
}

func (n *GrantRoleStmt) Pos() int {
	return 0
}
