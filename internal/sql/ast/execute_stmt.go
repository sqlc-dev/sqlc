package ast

import ()

type ExecuteStmt struct {
	Name   *string
	Params *List
}

func (n *ExecuteStmt) Pos() int {
	return 0
}
