package ast

type SelectStmt struct {
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

func (n *SelectStmt) Pos() int {
	return 0
}

func set(n Node) bool {
	if n == nil {
		return false
	}
	_, ok := n.(*TODO)
	if ok {
		return false
	}
	return true
}

func items(n *List) bool {
	if n == nil {
		return false
	}
	return len(n.Items) > 0
}

func (n *SelectStmt) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}

	if n.WithClause != nil {
		buf.astFormat(n.WithClause)
		buf.WriteString(" ")
	}

	buf.WriteString("SELECT ")
	buf.astFormat(n.TargetList)

	if items(n.FromClause) {
		buf.WriteString(" FROM ")
		buf.astFormat(n.FromClause)
	}

	if set(n.WhereClause) {
		buf.WriteString(" WHERE ")
		buf.astFormat(n.WhereClause)
	}
}
