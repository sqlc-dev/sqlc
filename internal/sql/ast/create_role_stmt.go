package ast

type CreateRoleStmt struct {
	Tag NodeTag[CreateRoleStmt] `json:"tag"`

	StmtType RoleStmtType
	Role     *string `json:",omitempty"`
	Options  *List   `json:",omitempty"`
}

func (n *CreateRoleStmt) Pos() int {
	return 0
}
