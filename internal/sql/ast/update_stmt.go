package ast

type UpdateStmt struct {
	Relations     *List
	TargetList    *List
	WhereClause   Node
	FromClause    *List
	ReturningList *List
	WithClause    *WithClause
}

func (n *UpdateStmt) Pos() int {
	return 0
}
