package ast

import (
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
			buf.group()
			buf.indent()
			buf.softline()
			if r, ok := row.(*List); ok {
				buf.joinComma(r, d)
			} else {
				buf.astFormat(row, d)
			}
			buf.endIndent()
			buf.softline()
			buf.endGroup()
			buf.WriteString(")")
		}
		return
	}

	buf.group()
	defer buf.endGroup()

	if n.WithClause != nil {
		buf.astFormat(n.WithClause, d)
		buf.line()
	}

	if n.Larg != nil && n.Rarg != nil {
		buf.astFormat(n.Larg, d)
		buf.line()
		switch n.Op {
		case Union:
			buf.WriteString("UNION")
		case Except:
			buf.WriteString("EXCEPT")
		case Intersect:
			buf.WriteString("INTERSECT")
		}
		if n.All {
			buf.WriteString(" ALL")
		}
		buf.line()
		buf.astFormat(n.Rarg, d)
	} else {
		buf.WriteString("SELECT")
		if items(n.DistinctClause) {
			buf.WriteString(" DISTINCT")
			if !todo(n.DistinctClause) {
				buf.WriteString(" ON (")
				buf.astFormat(n.DistinctClause, d)
				buf.WriteString(")")
			}
		}
		buf.group()
		buf.indent()
		buf.line()
		buf.joinComma(n.TargetList, d)
		buf.endIndent()
		buf.endGroup()
	}

	if items(n.FromClause) {
		buf.beforeClause(n.FromClause, d)
		buf.line()
		buf.WriteString("FROM ")
		buf.astFormat(n.FromClause, d)
	}

	if set(n.WhereClause) {
		buf.beforeClause(n.WhereClause, d)
		buf.line()
		buf.WriteString("WHERE ")
		buf.condition(n.WhereClause, d)
	}

	if items(n.GroupClause) {
		buf.beforeClause(n.GroupClause, d)
		buf.line()
		buf.WriteString("GROUP BY ")
		buf.astFormat(n.GroupClause, d)
	}

	if set(n.HavingClause) {
		buf.beforeClause(n.HavingClause, d)
		buf.line()
		buf.WriteString("HAVING ")
		buf.condition(n.HavingClause, d)
	}

	if items(n.SortClause) {
		buf.beforeClause(n.SortClause, d)
		buf.line()
		buf.WriteString("ORDER BY ")
		buf.astFormat(n.SortClause, d)
	}

	if set(n.LimitCount) {
		buf.beforeClause(n.LimitCount, d)
		buf.line()
		buf.WriteString("LIMIT ")
		buf.astFormat(n.LimitCount, d)
	}

	if set(n.LimitOffset) {
		buf.beforeClause(n.LimitOffset, d)
		buf.line()
		buf.WriteString("OFFSET ")
		buf.astFormat(n.LimitOffset, d)
	}

	if items(n.LockingClause) {
		buf.line()
		buf.astFormat(n.LockingClause, d)
	}

}
