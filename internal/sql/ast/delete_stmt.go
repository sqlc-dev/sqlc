package ast

type DeleteStmt struct {
	Relations     *List
	UsingClause   *List
	WhereClause   Node
	LimitCount    Node
	ReturningList *List
	WithClause    *WithClause
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

	buf.WriteString("DELETE FROM ")
	if items(n.Relations) {
		buf.astFormat(n.Relations)
	}

	if set(n.WhereClause) {
		buf.WriteString(" WHERE ")
		buf.astFormat(n.WhereClause)
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
