package ast

type AlterRoleSetStmt struct {
	Tag NodeTag[AlterRoleSetStmt] `json:"tag"`

	Role     *RoleSpec        `json:"role,omitempty"`
	Database *string          `json:"database,omitempty"`
	Setstmt  *VariableSetStmt `json:"setstmt,omitempty"`
}

func (n *AlterRoleSetStmt) Pos() int {
	return 0
}
