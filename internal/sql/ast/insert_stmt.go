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
	if n.OnConflictClause != nil {
		switch n.OnConflictClause.Action {
		case OnConflictAction_INSERT_OR_ABORT:
			buf.WriteString("INSERT OR ABORT INTO ")
		case OnConflictAction_INSERT_OR_REVERT:
			buf.WriteString("INSERT OR REVERT INTO ")
		case OnConflictAction_INSERT_OR_IGNORE:
			buf.WriteString("INSERT OR IGNORE INTO ")
		case OnConflictAction_UPSERT:
			buf.WriteString("UPSERT INTO ")
		case OnConflictAction_REPLACE:
			buf.WriteString("REPLACE INTO ")
		default:
			buf.WriteString("INSERT INTO ")
		}
	} else {
		buf.WriteString("INSERT INTO ")
	}
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

	if n.OnConflictClause != nil && n.OnConflictClause.Action < 4 {
		buf.WriteString(" ON CONFLICT DO NOTHING ")
	}

	if items(n.ReturningList) {
		buf.WriteString(" RETURNING ")
		buf.astFormat(n.ReturningList)
	}
}
