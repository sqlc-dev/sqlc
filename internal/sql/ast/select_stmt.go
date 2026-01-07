package ast

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/sql/format"
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

func (n *SelectStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}

	if items(n.ValuesLists) {
		buf.WriteString("VALUES ")
		// ValuesLists is a list of rows, where each row is a List of values
		for i, row := range n.ValuesLists.Items {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString("(")
			buf.astFormat(row, d)
			buf.WriteString(")")
		}
		return
	}

	if n.WithClause != nil {
		buf.astFormat(n.WithClause, d)
		buf.WriteString(" ")
	}

	if n.Larg != nil && n.Rarg != nil {
		buf.astFormat(n.Larg, d)
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
		buf.astFormat(n.Rarg, d)
	} else {
		buf.WriteString("SELECT ")
	}

	if items(n.DistinctClause) {
		buf.WriteString("DISTINCT ")
		if !todo(n.DistinctClause) {
			fmt.Fprintf(buf, "ON (")
			buf.astFormat(n.DistinctClause, d)
			fmt.Fprintf(buf, ")")
		}
	}
	buf.astFormat(n.TargetList, d)

	if items(n.FromClause) {
		buf.WriteString(" FROM ")
		buf.astFormat(n.FromClause, d)
	}

	if set(n.WhereClause) {
		buf.WriteString(" WHERE ")
		buf.astFormat(n.WhereClause, d)
	}

	if items(n.GroupClause) {
		buf.WriteString(" GROUP BY ")
		buf.astFormat(n.GroupClause, d)
	}

	if set(n.HavingClause) {
		buf.WriteString(" HAVING ")
		buf.astFormat(n.HavingClause, d)
	}

	if items(n.SortClause) {
		buf.WriteString(" ORDER BY ")
		buf.astFormat(n.SortClause, d)
	}

	if set(n.LimitCount) {
		buf.WriteString(" LIMIT ")
		buf.astFormat(n.LimitCount, d)
	}

	if set(n.LimitOffset) {
		buf.WriteString(" OFFSET ")
		buf.astFormat(n.LimitOffset, d)
	}

	if items(n.LockingClause) {
		buf.WriteString(" ")
		buf.astFormat(n.LockingClause, d)
	}

}
