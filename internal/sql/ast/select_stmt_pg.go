package ast

type SelectStmt_PG struct {
	DistinctClause *List
	IntoClause     *IntoClause
	TargetList     *List
	FromClause     *List
	WhereClause    Node
	GroupClause    *List
	HavingClause   Node
	WindowClause   *List
	ValuesLists    *List
	SortClause     *List
	LimitOffset    Node
	LimitCount     Node
	LockingClause  *List
	WithClause     *WithClause
	Op             SetOperation
	All            bool
	Larg           *SelectStmt
	Rarg           *SelectStmt
}

func (n *SelectStmt_PG) Pos() int {
	return 0
}
