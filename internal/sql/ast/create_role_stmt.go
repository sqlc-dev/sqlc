package ast

type CreateRoleStmt struct {
	StmtType RoleStmtType
	Role     *string
	Options  *List
}

func (n *CreateRoleStmt) Pos() int {
	return 0
}
