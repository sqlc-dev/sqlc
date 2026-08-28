package ast

type AlterRoleStmt struct {
	Tag NodeTag[AlterRoleStmt] `json:"tag"`

	Role    *RoleSpec
	Options *List
	Action  int
}

func (n *AlterRoleStmt) Pos() int {
	return 0
}
