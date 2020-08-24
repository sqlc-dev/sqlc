package ast

import ()

type AlterTSDictionaryStmt struct {
	Dictname *List
	Options  *List
}

func (n *AlterTSDictionaryStmt) Pos() int {
	return 0
}
