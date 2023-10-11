package ast

import (
	"fmt"
)

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

func (n *SelectStmt) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}

	if items(n.ValuesLists) {
		buf.WriteString("VALUES (")
		buf.astFormat(n.ValuesLists)
		buf.WriteString(")")
		return
	}

	if n.WithClause != nil {
		buf.astFormat(n.WithClause)
		buf.WriteString(" ")
	}

	if n.Larg != nil && n.Rarg != nil {
		buf.astFormat(n.Larg)
		switch n.Op {
		case Union:
			buf.WriteString(" UNION ")
		case Except:
			buf.WriteString(" EXCEPT ")
		case Intersect:
			buf.WriteString(" INTERSECT ")
		}
		if n.All {
			buf.WriteString("ALL ")
		}
		buf.astFormat(n.Rarg)
	} else {
		buf.WriteString("SELECT ")
	}

	if items(n.DistinctClause) {
		buf.WriteString("DISTINCT ")
		if !todo(n.DistinctClause) {
			fmt.Fprintf(buf, "ON (")
			buf.astFormat(n.DistinctClause)
			fmt.Fprintf(buf, ")")
		}
	}
	buf.astFormat(n.TargetList)

	if items(n.FromClause) {
		buf.WriteString(" FROM ")
		buf.astFormat(n.FromClause)
	}

	if set(n.WhereClause) {
		buf.WriteString(" WHERE ")
		buf.astFormat(n.WhereClause)
	}

	if items(n.GroupClause) {
		buf.WriteString(" GROUP BY ")
		buf.astFormat(n.GroupClause)
	}

	if items(n.SortClause) {
		buf.WriteString(" ORDER BY ")
		buf.astFormat(n.SortClause)
	}

	if set(n.LimitCount) {
		buf.WriteString(" LIMIT ")
		buf.astFormat(n.LimitCount)
	}

	if set(n.LimitOffset) {
		buf.WriteString(" OFFSET ")
		buf.astFormat(n.LimitOffset)
	}

	if items(n.LockingClause) {
		buf.WriteString(" ")
		buf.astFormat(n.LockingClause)
	}

}
