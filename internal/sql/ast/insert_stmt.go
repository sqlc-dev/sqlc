package ast

type InsertStmt struct {
	Relation         *RangeVar
	Cols             *List
	SelectStmt       Node
	OnConflictClause *OnConflictClause
	ReturningList    *List
	WithClause       *WithClause
	Override         OverridingKind
}

func (n *InsertStmt) Pos() int {
	return 0
}

func (n *InsertStmt) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}

	if n.WithClause != nil {
		buf.astFormat(n.WithClause)
		buf.WriteString(" ")
	}

	buf.WriteString("INSERT INTO ")
	if n.Relation != nil {
		buf.astFormat(n.Relation)
	}
	if items(n.Cols) {
		buf.WriteString(" (")
		buf.astFormat(n.Cols)
		buf.WriteString(") ")
	}

	if set(n.SelectStmt) {
		buf.astFormat(n.SelectStmt)
	}

	if n.OnConflictClause != nil {
		buf.WriteString(" ON CONFLICT DO NOTHING ")
	}

	if items(n.ReturningList) {
		buf.WriteString(" RETURNING ")
		buf.astFormat(n.ReturningList)
	}
}
