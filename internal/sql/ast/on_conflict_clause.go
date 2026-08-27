package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type OnConflictClause struct {
	Action      OnConflictAction
	Infer       *InferClause
	TargetList  *List
	WhereClause Node
	Location    int
}

func (n *OnConflictClause) Pos() int {
	return n.Location
}

// OnConflictAction values matching pg_query_go
const (
	OnConflictActionUndefined OnConflictAction = 0
	OnConflictActionNone      OnConflictAction = 1
	OnConflictActionNothing   OnConflictAction = 2
	OnConflictActionUpdate    OnConflictAction = 3
)

func (n *OnConflictClause) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("ON CONFLICT ")
	if n.Infer != nil {
		buf.astFormat(n.Infer, d)
		buf.WriteString(" ")
	}
	switch n.Action {
	case OnConflictActionNothing:
		buf.WriteString("DO NOTHING")
	case OnConflictActionUpdate:
		buf.WriteString("DO UPDATE SET")
		buf.Group()
		buf.Indent()
		// Format as assignment list: name = val, one per line when the
		// author broke the list (or a comment forces it open).
		if n.TargetList != nil {
			for i, item := range n.TargetList.Items {
				if i == 0 {
					buf.boundary(item)
					buf.Line()
				} else {
					buf.WriteString(",")
					buf.boundary(item)
					buf.Line()
				}
				if rt, ok := item.(*ResTarget); ok {
					if rt.Name != nil {
						buf.WriteString(*rt.Name)
					}
					buf.WriteString(" = ")
					buf.astFormat(rt.Val, d)
				} else {
					buf.astFormat(item, d)
				}
			}
		}
		buf.EndIndent()
		if set(n.WhereClause) {
			buf.boundary(n.WhereClause)
			buf.Line()
			buf.WriteString("WHERE ")
			buf.condition(n.WhereClause, d)
		}
		buf.EndGroup()
	}
}
