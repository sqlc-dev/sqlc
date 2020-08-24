package ast

import ()

type CreateEnumStmt struct {
	TypeName *List
	Vals     *List
}

func (n *CreateEnumStmt) Pos() int {
	return 0
}
