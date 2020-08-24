package ast

import ()

type DropRoleStmt struct {
	Roles     *List
	MissingOk bool
}

func (n *DropRoleStmt) Pos() int {
	return 0
}
