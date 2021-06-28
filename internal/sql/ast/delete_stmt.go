package ast

type DeleteStmt struct {
	Relation      *RangeVar
	UsingClause   *List
	WhereClause   Node
	ReturningList *List
	WithClause    *WithClause
}

func (n *DeleteStmt) Pos() int {
	return 0
}
