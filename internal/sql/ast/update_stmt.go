package ast

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/sql/format"
)

type UpdateStmt struct {
	Tag NodeTag[UpdateStmt] `json:"tag"`

	Relations     *List       `json:"relations,omitempty"`
	TargetList    *List       `json:"target_list,omitempty"`
	WhereClause   Node        `json:"where_clause,omitempty"`
	FromClause    *List       `json:"from_clause,omitempty"`
	LimitCount    Node        `json:"limit_count,omitempty"`
	ReturningList *List       `json:"returning_list,omitempty"`
	WithClause    *WithClause `json:"with_clause,omitempty"`
	// PostgreSQL 18 RETURNING WITH (OLD AS ..., NEW AS ...) aliases
	ReturningOldAlias string `json:"returning_old_alias"`
	ReturningNewAlias string `json:"returning_new_alias"`
	// TableRefs is the statement's table-reference tree as the author wrote
	// it (MySQL: UPDATE a JOIN b ON ...). Relations flattens that tree to
	// bare tables for analysis; printing prefers TableRefs when it is set,
	// so a join's ON condition survives formatting.
	TableRefs *List `json:"table_refs,omitempty"`
}

func (n *UpdateStmt) Pos() int {
	return 0
}

func (n *UpdateStmt) Format(buf *TrackedBuffer, d format.Dialect) {
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

	buf.WriteString("UPDATE ")
	if items(n.TableRefs) {
		buf.astFormat(n.TableRefs, d)
	} else if items(n.Relations) {
		buf.astFormat(n.Relations, d)
	}

	if items(n.TargetList) {
		buf.beforeClause(n.TargetList, d)
		buf.Line()
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
			buf.Group()
			buf.Indent()
			for i, item := range n.TargetList.Items {
				if i > 0 {
					buf.WriteString(",")
					buf.boundary(item)
					buf.Line()
				}
				switch nn := item.(type) {
				case *ResTarget:
					if nn.Relation != nil {
						buf.WriteString(d.QuoteIdent(*nn.Relation))
						buf.WriteString(".")
					}
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
			buf.EndIndent()
			buf.EndGroup()
		}
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
