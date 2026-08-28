package ast

type AlterRoleSetStmt struct {
	Tag NodeTag[AlterRoleSetStmt] `json:"tag"`

	Role     *RoleSpec
	Database *string
	Setstmt  *VariableSetStmt
}

func (n *AlterRoleSetStmt) Pos() int {
	return 0
}
