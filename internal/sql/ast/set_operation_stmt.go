package ast

type SetOperationStmt struct {
	Op            SetOperation
	All           bool
	Larg          Node
	Rarg          Node
	ColTypes      *List
	ColTypmods    *List
	ColCollations *List
	GroupClauses  *List
}

func (n *SetOperationStmt) Pos() int {
	return 0
}
