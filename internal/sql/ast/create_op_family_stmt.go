package ast

import ()

type CreateOpFamilyStmt struct {
	Opfamilyname *List
	Amname       *string
}

func (n *CreateOpFamilyStmt) Pos() int {
	return 0
}
