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
	// PostgreSQL 18 RETURNING WITH (OLD AS ..., NEW AS ...) aliases
	ReturningOldAlias string
	ReturningNewAlias string
}

func (n *DeleteStmt) Pos() int {
	return 0
}

func (n *DeleteStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}

	buf.group()
	defer buf.endGroup()

	if n.WithClause != nil {
		buf.astFormat(n.WithClause, d)
		buf.line()
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
		buf.line()
		buf.WriteString("USING ")
		buf.join(n.UsingClause, d, ", ")
	}

	if set(n.WhereClause) {
		buf.line()
		buf.WriteString("WHERE ")
		buf.condition(n.WhereClause, d)
	}

	if set(n.LimitCount) {
		buf.line()
		buf.WriteString("LIMIT ")
		buf.astFormat(n.LimitCount, d)
	}

	if items(n.ReturningList) {
		buf.line()
		buf.WriteString("RETURNING ")
		formatReturningOptions(buf, d, n.ReturningOldAlias, n.ReturningNewAlias)
		buf.astFormat(n.ReturningList, d)
	}
}
