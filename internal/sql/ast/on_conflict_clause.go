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
		buf.WriteString("DO UPDATE SET ")
		formatSetClause(buf, d, n.TargetList)
		if set(n.WhereClause) {
			buf.WriteString(" WHERE ")
			buf.astFormat(n.WhereClause, d)
		}
	}
}
