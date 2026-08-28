package ast

type CreateRoleStmt struct {
	Tag NodeTag[CreateRoleStmt] `json:"tag"`

	StmtType RoleStmtType `json:"stmt_type"`
	Role     *string      `json:"role,omitempty"`
	Options  *List        `json:"options,omitempty"`
}

func (n *CreateRoleStmt) Pos() int {
	return 0
}
