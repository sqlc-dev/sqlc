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
	// PostgreSQL 18 RETURNING WITH (OLD AS ..., NEW AS ...) aliases
	ReturningOldAlias string
	ReturningNewAlias string
}

func (n *UpdateStmt) Pos() int {
	return 0
}

func (n *UpdateStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.group()
	defer buf.endGroup()

	if n.WithClause != nil {
		buf.astFormat(n.WithClause, d)
		buf.line()
	}

	buf.WriteString("UPDATE ")
	if items(n.Relations) {
		buf.astFormat(n.Relations, d)
	}

	if items(n.TargetList) {
		buf.beforeClause(n.TargetList, d)
		buf.line()
		buf.WriteString("SET ")

		multi := false
		for _, item := range n.TargetList.Items {
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
			for _, item := range n.TargetList.Items {
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
			buf.group()
			buf.indent()
			for i, item := range n.TargetList.Items {
				if i > 0 {
					buf.WriteString(",")
					buf.line()
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
			buf.endIndent()
			buf.endGroup()
		}
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

	if set(n.LimitCount) {
		buf.beforeClause(n.LimitCount, d)
		buf.line()
		buf.WriteString("LIMIT ")
		buf.astFormat(n.LimitCount, d)
	}

	if items(n.ReturningList) {
		buf.beforeClause(n.ReturningList, d)
		buf.line()
		buf.WriteString("RETURNING ")
		formatReturningOptions(buf, d, n.ReturningOldAlias, n.ReturningNewAlias)
		buf.astFormat(n.ReturningList, d)
	}
}
