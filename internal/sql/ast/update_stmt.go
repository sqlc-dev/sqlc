package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type UpdateStmt struct {
	Relations     *List
	TargetList    *List
	WhereClause   Node
	FromClause    *List
	LimitCount    Node
	ReturningList *List
	WithClause    *WithClause
	// PostgreSQL 18 RETURNING WITH (OLD AS ..., NEW AS ...) aliases
	ReturningOldAlias string
	ReturningNewAlias string
}

func (n *UpdateStmt) Pos() int {
	return 0
}

// formatSetClause renders an UPDATE-style SET assignment list (without the
// leading "SET ") for the given target list. It handles both single-column
// assignments (col = val) and multi-column assignments ((a, b) = (x, y)), and
// is shared by UPDATE and MERGE ... WHEN MATCHED THEN UPDATE.
func formatSetClause(buf *TrackedBuffer, d format.Dialect, targetList *List) {
	if !items(targetList) {
		return
	}

	for i := 0; i < len(targetList.Items); {
		if i > 0 {
			buf.WriteString(", ")
		}
		target, ok := targetList.Items[i].(*ResTarget)
		if !ok {
			buf.astFormat(targetList.Items[i], d)
			i++
			continue
		}

		multi, ok := target.Val.(*MultiAssignRef)
		if !ok {
			formatSetTarget(buf, d, target)
			buf.WriteString(" = ")
			buf.astFormat(target.Val, d)
			i++
			continue
		}

		count := multi.Ncolumns
		if count < 1 {
			count = 1
		}
		end := min(i+count, len(targetList.Items))
		buf.WriteString("(")
		for j := i; j < end; j++ {
			if j > i {
				buf.WriteString(", ")
			}
			if grouped, ok := targetList.Items[j].(*ResTarget); ok {
				formatSetTarget(buf, d, grouped)
			} else {
				buf.astFormat(targetList.Items[j], d)
			}
		}
		buf.WriteString(") = ")
		formatMultiAssignSource(buf, d, multi.Source)
		i = end
	}
}

func formatSetTarget(buf *TrackedBuffer, d format.Dialect, target *ResTarget) {
	if target.Name != nil {
		buf.WriteString(d.QuoteIdent(*target.Name))
	}
	// Handle array subscript indirection (e.g., names[$1]).
	if items(target.Indirection) {
		for _, ind := range target.Indirection.Items {
			buf.astFormat(ind, d)
		}
	}
}

func formatMultiAssignSource(buf *TrackedBuffer, d format.Dialect, source Node) {
	if row, ok := source.(*RowExpr); ok {
		buf.WriteString("(")
		if items(row.Args) {
			buf.join(row.Args, d, ", ")
		}
		buf.WriteString(")")
		return
	}
	buf.astFormat(source, d)
}

func (n *UpdateStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	if n.WithClause != nil {
		buf.astFormat(n.WithClause, d)
		buf.WriteString(" ")
	}

	buf.WriteString("UPDATE ")
	if items(n.Relations) {
		buf.astFormat(n.Relations, d)
	}

	if items(n.TargetList) {
		buf.WriteString(" SET ")
		formatSetClause(buf, d, n.TargetList)
	}

	if items(n.FromClause) {
		buf.WriteString(" FROM ")
		buf.astFormat(n.FromClause, d)
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
		formatReturningOptions(buf, d, n.ReturningOldAlias, n.ReturningNewAlias)
		buf.astFormat(n.ReturningList, d)
	}
}
