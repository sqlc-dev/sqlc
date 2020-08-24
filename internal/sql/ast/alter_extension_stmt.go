package ast

import ()

type AlterExtensionStmt struct {
	Extname *string
	Options *List
}

func (n *AlterExtensionStmt) Pos() int {
	return 0
}
