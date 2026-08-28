package ast

type CreateRoleStmt struct {
	Tag NodeTag[CreateRoleStmt] `json:"tag"`

	StmtType RoleStmtType
	Role     *string
	Options  *List
}

func (n *CreateRoleStmt) Pos() int {
	return 0
}
