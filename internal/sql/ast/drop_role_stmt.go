package ast

type DropRoleStmt struct {
	Tag NodeTag[DropRoleStmt] `json:"tag"`

	Roles     *List `json:"roles,omitempty"`
	MissingOk bool  `json:"missing_ok"`
}

func (n *DropRoleStmt) Pos() int {
	return 0
}
