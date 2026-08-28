package ast

type AlterRoleStmt struct {
	Tag NodeTag[AlterRoleStmt] `json:"tag"`

	Role    *RoleSpec `json:",omitempty"`
	Options *List     `json:",omitempty"`
	Action  int
}

func (n *AlterRoleStmt) Pos() int {
	return 0
}
