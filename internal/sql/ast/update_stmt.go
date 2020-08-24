package ast

import ()

type UpdateStmt struct {
	Relation      *RangeVar
	TargetList    *List
	WhereClause   Node
	FromClause    *List
	ReturningList *List
	WithClause    *WithClause
}

func (n *UpdateStmt) Pos() int {
	return 0
}
