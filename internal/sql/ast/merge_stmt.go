package ast

type MergeStmt struct {
	Relation         *RangeVar
	JoinCondition    Node
	WithClause       *WithClause
	SourceRelation   Node
	MergeWhenClauses *List
}

func (n *MergeStmt) Pos() int {
	return 0
}
