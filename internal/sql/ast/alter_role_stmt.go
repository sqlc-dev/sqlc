package ast

type AlterRoleStmt struct {
	Tag NodeTag[AlterRoleStmt] `json:"tag"`

	Role    *RoleSpec `json:"role,omitempty"`
	Options *List     `json:"options,omitempty"`
	Action  int       `json:"action"`
}

func (n *AlterRoleStmt) Pos() int {
	return 0
}
