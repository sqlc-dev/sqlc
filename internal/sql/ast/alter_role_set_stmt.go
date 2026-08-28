package ast

type AlterRoleSetStmt struct {
	Tag NodeTag[AlterRoleSetStmt] `json:"tag"`

	Role     *RoleSpec        `json:",omitempty"`
	Database *string          `json:",omitempty"`
	Setstmt  *VariableSetStmt `json:",omitempty"`
}

func (n *AlterRoleSetStmt) Pos() int {
	return 0
}
