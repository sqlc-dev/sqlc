package ast

type CreateRoleStmt struct {
	StmtType RoleStmtType
	Role     *string
	Options  *List

	// YDB specific
	BindRole Node
}

func (n *CreateRoleStmt) Pos() int {
	return 0
}
