package ast

type DropRoleStmt struct {
	Tag NodeTag[DropRoleStmt] `json:"tag"`

	Roles     *List
	MissingOk bool
}

func (n *DropRoleStmt) Pos() int {
	return 0
}
