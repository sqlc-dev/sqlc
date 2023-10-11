package ast

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

func (n *UpdateStmt) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	if n.WithClause != nil {
		buf.astFormat(n.WithClause)
		buf.WriteString(" ")
	}

	buf.WriteString("UPDATE ")
	if items(n.Relations) {
		buf.astFormat(n.Relations)
	}

	if items(n.TargetList) {
		buf.WriteString(" SET (")
		buf.astFormat(n.TargetList)
		buf.WriteString(") ")
	}

	if set(n.WhereClause) {
		buf.WriteString(" WHERE ")
		buf.astFormat(n.WhereClause)
	}

	if set(n.LimitCount) {
		buf.WriteString(" LIMIT ")
		buf.astFormat(n.LimitCount)
	}
}
