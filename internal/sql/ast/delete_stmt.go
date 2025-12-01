package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type DeleteStmt struct {
	Relations     *List
	UsingClause   *List
	WhereClause   Node
	LimitCount    Node
	ReturningList *List
	WithClause    *WithClause
	// MySQL multi-table DELETE support
	Targets    *List // Tables to delete from (e.g., jt.*, pt.*)
	FromClause Node  // FROM clause with JOINs (Node to support JoinExpr)
}

func (n *DeleteStmt) Pos() int {
	return 0
}

func (n *DeleteStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}

	if n.WithClause != nil {
		buf.astFormat(n.WithClause, d)
		buf.WriteString(" ")
	}

	buf.WriteString("DELETE ")

	// MySQL multi-table DELETE: DELETE t1.*, t2.* FROM t1 JOIN t2 ...
	if items(n.Targets) {
		buf.join(n.Targets, d, ", ")
		buf.WriteString(" FROM ")
		if set(n.FromClause) {
			buf.astFormat(n.FromClause, d)
		} else if items(n.Relations) {
			buf.astFormat(n.Relations, d)
		}
	} else {
		buf.WriteString("FROM ")
		if items(n.Relations) {
			buf.astFormat(n.Relations, d)
		}
	}

	if items(n.UsingClause) {
		buf.WriteString(" USING ")
		buf.join(n.UsingClause, d, ", ")
	}

	if set(n.WhereClause) {
		buf.WriteString(" WHERE ")
		buf.astFormat(n.WhereClause, d)
	}

	if set(n.LimitCount) {
		buf.WriteString(" LIMIT ")
		buf.astFormat(n.LimitCount, d)
	}

	if items(n.ReturningList) {
		buf.WriteString(" RETURNING ")
		buf.astFormat(n.ReturningList, d)
	}
}
