package ast

import (
	"github.com/sqlc-dev/sqlc/internal/sql/format"
)

type SelectStmt struct {
	Tag NodeTag[SelectStmt] `json:"tag"`

	DistinctClause *List        `json:"distinct_clause,omitempty"`
	IntoClause     *IntoClause  `json:"into_clause,omitempty"`
	TargetList     *List        `json:"target_list,omitempty"`
	FromClause     *List        `json:"from_clause,omitempty"`
	WhereClause    Node         `json:"where_clause,omitempty"`
	GroupClause    *List        `json:"group_clause,omitempty"`
	HavingClause   Node         `json:"having_clause,omitempty"`
	WindowClause   *List        `json:"window_clause,omitempty"`
	ValuesLists    *List        `json:"values_lists,omitempty"`
	SortClause     *List        `json:"sort_clause,omitempty"`
	LimitOffset    Node         `json:"limit_offset,omitempty"`
	LimitCount     Node         `json:"limit_count,omitempty"`
	LockingClause  *List        `json:"locking_clause,omitempty"`
	WithClause     *WithClause  `json:"with_clause,omitempty"`
	Op             SetOperation `json:"op"`
	All            bool         `json:"all"`
	Larg           *SelectStmt  `json:"larg,omitempty"`
	Rarg           *SelectStmt  `json:"rarg,omitempty"`
	// TableHints is the text inside a MySQL optimizer-hint comment
	// (SELECT /*+ MAX_EXECUTION_TIME(1000) */ ...). The compiler ignores
	// hints; printing keeps them because they change how the server runs
	// the query.
	TableHints string `json:"table_hints,omitempty"`
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
		if n.TableHints != "" {
			buf.WriteString(" /*+ ")
			buf.WriteString(n.TableHints)
			buf.WriteString(" */")
		}
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
