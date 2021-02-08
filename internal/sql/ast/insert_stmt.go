package ast

type InsertStmt struct {
	Relation         *RangeVar
	Cols             *List
	SelectStmt       Node
	OnConflictClause *OnConflictClause
	ReturningList    *List
	WithClause       *WithClause
	Override         OverridingKind
}

func (n *InsertStmt) Pos() int {
	return 0
}
