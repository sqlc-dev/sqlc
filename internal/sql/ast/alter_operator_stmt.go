package ast

import ()

type AlterOperatorStmt struct {
	Opername *ObjectWithArgs
	Options  *List
}

func (n *AlterOperatorStmt) Pos() int {
	return 0
}
