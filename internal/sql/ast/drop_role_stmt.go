package ast

type DropRoleStmt struct {
	Tag NodeTag[DropRoleStmt] `json:"tag"`

	Roles     *List `json:",omitempty"`
	MissingOk bool
}

func (n *DropRoleStmt) Pos() int {
	return 0
}
