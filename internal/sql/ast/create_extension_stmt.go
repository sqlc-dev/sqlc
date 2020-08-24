package ast

import ()

type CreateExtensionStmt struct {
	Extname     *string
	IfNotExists bool
	Options     *List
}

func (n *CreateExtensionStmt) Pos() int {
	return 0
}
