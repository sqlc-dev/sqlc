package ast

import (
	"github.com/sqlc-dev/sqlc/internal/sql/format"
)

type SelectStmt struct {
	Tag NodeTag[SelectStmt] `json:"tag"`

	DistinctClause *List       `json:",omitempty"`
	IntoClause     *IntoClause `json:",omitempty"`
	TargetList     *List       `json:",omitempty"`
	FromClause     *List       `json:",omitempty"`
	WhereClause    Node        `json:",omitempty"`
	GroupClause    *List       `json:",omitempty"`
	HavingClause   Node        `json:",omitempty"`
	WindowClause   *List       `json:",omitempty"`
	ValuesLists    *List       `json:",omitempty"`
	SortClause     *List       `json:",omitempty"`
	LimitOffset    Node        `json:",omitempty"`
	LimitCount     Node        `json:",omitempty"`
	LockingClause  *List       `json:",omitempty"`
	WithClause     *WithClause `json:",omitempty"`
	Op             SetOperation
	All            bool
	Larg           *SelectStmt `json:",omitempty"`
	Rarg           *SelectStmt `json:",omitempty"`
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
			buf.Group()
			buf.Indent()
			buf.Softline()
			if r, ok := row.(*List); ok {
				buf.joinComma(r, d)
			} else {
				buf.astFormat(row, d)
			}
			buf.EndIndent()
			buf.Softline()
			buf.EndGroup()
			buf.WriteString(")")
		}
		return
	}

	buf.Group()
	defer buf.EndGroup()

	if n.WithClause != nil {
		buf.astFormat(n.WithClause, d)
		// The boundary between the WITH clause and the statement's own
		// first keyword; the group keeps a break inside the WITH clause
		// from forcing the statement body apart clause by clause.
		if n.Larg != nil {
			buf.boundary(n.Larg)
		} else {
			buf.boundary(n.TargetList)
		}
		buf.Line()
		buf.Group()
		defer buf.EndGroup()
	}

	if n.Larg != nil && n.Rarg != nil {
		buf.astFormat(n.Larg, d)
		// The seam between the compound halves: an author who broke the
		// line around UNION / INTERSECT / EXCEPT keeps the operator on its
		// own line, and a comment above it prints here.
		buf.boundary(n.Rarg)
		buf.Line()
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
		buf.Line()
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
		buf.Group()
		buf.Indent()
		buf.Line()
		buf.joinComma(n.TargetList, d)
		buf.EndIndent()
		buf.EndGroup()
	}

	if items(n.FromClause) {
		buf.beforeClause(n.FromClause, d)
		buf.Line()
		buf.WriteString("FROM ")
		buf.astFormat(n.FromClause, d)
	}

	if set(n.WhereClause) {
		buf.beforeClause(n.WhereClause, d)
		buf.Line()
		buf.WriteString("WHERE ")
		buf.condition(n.WhereClause, d)
	}

	if items(n.GroupClause) {
		buf.beforeClause(n.GroupClause, d)
		buf.Line()
		buf.WriteString("GROUP BY ")
		buf.astFormat(n.GroupClause, d)
	}

	if set(n.HavingClause) {
		buf.beforeClause(n.HavingClause, d)
		buf.Line()
		buf.WriteString("HAVING ")
		buf.condition(n.HavingClause, d)
	}

	if items(n.SortClause) {
		buf.beforeClause(n.SortClause, d)
		buf.Line()
		buf.WriteString("ORDER BY ")
		buf.astFormat(n.SortClause, d)
	}

	if set(n.LimitCount) {
		buf.beforeClause(n.LimitCount, d)
		buf.Line()
		buf.WriteString("LIMIT ")
		buf.astFormat(n.LimitCount, d)
	}

	if set(n.LimitOffset) {
		buf.beforeClause(n.LimitOffset, d)
		buf.Line()
		buf.WriteString("OFFSET ")
		buf.astFormat(n.LimitOffset, d)
	}

	if items(n.LockingClause) {
		buf.Line()
		buf.astFormat(n.LockingClause, d)
	}

}
