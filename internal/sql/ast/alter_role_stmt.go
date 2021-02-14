package ast

type AlterRoleStmt struct {
	Role    *RoleSpec
	Options *List
	Action  int
}

func (n *AlterRoleStmt) Pos() int {
	return 0
}
