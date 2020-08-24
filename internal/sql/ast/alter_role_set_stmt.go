package pg

type AlterRoleSetStmt struct {
	Role     *RoleSpec
	Database *string
	Setstmt  *VariableSetStmt
}

func (n *AlterRoleSetStmt) Pos() int {
	return 0
}
