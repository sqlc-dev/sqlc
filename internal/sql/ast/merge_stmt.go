package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type MergeStmt struct {
	Relation         *RangeVar
	SourceRelation   Node
	JoinCondition    Node
	MergeWhenClauses *List
	ReturningList    *List
	WithClause       *WithClause
}

func (n *MergeStmt) Pos() int {
	return 0
}

func (n *MergeStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}

	if n.WithClause != nil {
		buf.astFormat(n.WithClause, d)
		buf.WriteString(" ")
	}

	buf.WriteString("MERGE INTO ")
	if n.Relation != nil {
		buf.astFormat(n.Relation, d)
	}

	if set(n.SourceRelation) {
		buf.WriteString(" USING ")
		buf.astFormat(n.SourceRelation, d)
	}

	if set(n.JoinCondition) {
		buf.WriteString(" ON ")
		buf.astFormat(n.JoinCondition, d)
	}

	if n.MergeWhenClauses != nil {
		for _, item := range n.MergeWhenClauses.Items {
			buf.WriteString(" ")
			buf.astFormat(item, d)
		}
	}

	if items(n.ReturningList) {
		buf.WriteString(" RETURNING ")
		buf.astFormat(n.ReturningList, d)
	}
}
