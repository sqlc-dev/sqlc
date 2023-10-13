package ast

type NotifyStmt struct {
	Conditionname *string
	Payload       *string
}

func (n *NotifyStmt) Pos() int {
	return 0
}

func (n *NotifyStmt) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString("NOTIFY ")
	if n.Conditionname != nil {
		buf.WriteString(*n.Conditionname)
	}
	if n.Payload != nil {
		buf.WriteString(", '")
		buf.WriteString(*n.Payload)
		buf.WriteString("'")
	}
}
