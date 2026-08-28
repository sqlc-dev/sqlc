package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type DeleteStmt struct {
	Tag NodeTag[DeleteStmt] `json:"tag"`

	Relations     *List       `json:"relations,omitempty"`
	UsingClause   *List       `json:"using_clause,omitempty"`
	WhereClause   Node        `json:"where_clause,omitempty"`
	LimitCount    Node        `json:"limit_count,omitempty"`
	ReturningList *List       `json:"returning_list,omitempty"`
	WithClause    *WithClause `json:"with_clause,omitempty"`
	// MySQL multi-table DELETE support
	Targets    *List `json:"targets,omitempty"`     // Tables to delete from (e.g., jt.*, pt.*)
	FromClause Node  `json:"from_clause,omitempty"` // FROM clause with JOINs (Node to support JoinExpr)
	// PostgreSQL 18 RETURNING WITH (OLD AS ..., NEW AS ...) aliases
	ReturningOldAlias string `json:"returning_old_alias"`
	ReturningNewAlias string `json:"returning_new_alias"`
}

func (n *DeleteStmt) Pos() int {
	return 0
}

func (n *DeleteStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}

	buf.Group()
	defer buf.EndGroup()

	if n.WithClause != nil {
		buf.astFormat(n.WithClause, d)
		// The boundary between the WITH clause and the statement's own
		// first keyword; the group keeps a break inside the WITH clause
		// from forcing the statement body apart clause by clause.
		buf.boundary(n.Relations)
		buf.Line()
		buf.Group()
		defer buf.EndGroup()
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
		buf.beforeClause(n.UsingClause, d)
		buf.Line()
		buf.WriteString("USING ")
		buf.join(n.UsingClause, d, ", ")
	}

	if set(n.WhereClause) {
		buf.beforeClause(n.WhereClause, d)
		buf.Line()
		buf.WriteString("WHERE ")
		buf.condition(n.WhereClause, d)
	}

	if set(n.LimitCount) {
		buf.beforeClause(n.LimitCount, d)
		buf.Line()
		buf.WriteString("LIMIT ")
		buf.astFormat(n.LimitCount, d)
	}

	if items(n.ReturningList) {
		buf.beforeClause(n.ReturningList, d)
		buf.Line()
		buf.WriteString("RETURNING ")
		formatReturningOptions(buf, d, n.ReturningOldAlias, n.ReturningNewAlias)
		buf.astFormat(n.ReturningList, d)
	}
}
