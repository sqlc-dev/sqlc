package ast

type DeleteStmt struct {
	Relations   *List
	UsingClause *List
	WhereClause Node
	LimitCount  Node

	ReturningList *List
	WithClause    *WithClause

	// YDB specific
	Batch        bool
	OnCols       *List
	OnSelectStmt Node
}

func (n *DeleteStmt) Pos() int {
	return 0
}

func (n *DeleteStmt) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}

	if n.WithClause != nil {
		buf.astFormat(n.WithClause)
		buf.WriteString(" ")
	}
	if n.Batch {
		buf.WriteString("BATCH ")
	}

	buf.WriteString("DELETE FROM ")
	if items(n.Relations) {
		buf.astFormat(n.Relations)
	}

	if set(n.WhereClause) {
		buf.WriteString(" WHERE ")
		buf.astFormat(n.WhereClause)
	}

	if items(n.OnCols) || set(n.OnSelectStmt) {
		buf.WriteString(" ON ")

		if items(n.OnCols) {
			buf.WriteString("(")
			buf.astFormat(n.OnCols)
			buf.WriteString(") ")
		}

		if set(n.OnSelectStmt) {
			buf.astFormat(n.OnSelectStmt)
		}
	}

	if set(n.LimitCount) {
		buf.WriteString(" LIMIT ")
		buf.astFormat(n.LimitCount)
	}

	if items(n.ReturningList) {
		buf.WriteString(" RETURNING ")
		buf.astFormat(n.ReturningList)
	}
}
