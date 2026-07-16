package ast

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/sql/format"
)

type UpdateStmt struct {
	Relations     *List
	TargetList    *List
	WhereClause   Node
	FromClause    *List
	LimitCount    Node
	ReturningList *List
	WithClause    *WithClause
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

	multi := false
	for _, item := range targetList.Items {
		switch nn := item.(type) {
		case *ResTarget:
			if _, ok := nn.Val.(*MultiAssignRef); ok {
				multi = true
			}
		}
	}
	if multi {
		names := []string{}
		vals := &List{}
		for _, item := range targetList.Items {
			res, ok := item.(*ResTarget)
			if !ok {
				continue
			}
			if res.Name != nil {
				names = append(names, *res.Name)
			}
			multi, ok := res.Val.(*MultiAssignRef)
			if !ok {
				vals.Items = append(vals.Items, res.Val)
				continue
			}
			row, ok := multi.Source.(*RowExpr)
			if !ok {
				vals.Items = append(vals.Items, res.Val)
				continue
			}
			vals.Items = append(vals.Items, row.Args.Items[multi.Colno-1])
		}

		buf.WriteString("(")
		buf.WriteString(strings.Join(names, ","))
		buf.WriteString(") = (")
		buf.join(vals, d, ",")
		buf.WriteString(")")
	} else {
		for i, item := range targetList.Items {
			if i > 0 {
				buf.WriteString(", ")
			}
			switch nn := item.(type) {
			case *ResTarget:
				if nn.Name != nil {
					buf.WriteString(d.QuoteIdent(*nn.Name))
				}
				// Handle array subscript indirection (e.g., names[$1])
				if items(nn.Indirection) {
					for _, ind := range nn.Indirection.Items {
						buf.astFormat(ind, d)
					}
				}
				buf.WriteString(" = ")
				buf.astFormat(nn.Val, d)
			default:
				buf.astFormat(item, d)
			}
		}
	}
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
		buf.astFormat(n.ReturningList, d)
	}
}
