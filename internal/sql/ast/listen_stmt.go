package ast

type ListenStmt struct {
	Conditionname *string
}

func (n *ListenStmt) Pos() int {
	return 0
}

func (n *ListenStmt) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString("LISTEN ")
	if n.Conditionname != nil {
		buf.WriteString(*n.Conditionname)
	}
}
