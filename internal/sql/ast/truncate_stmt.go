package ast

type TruncateStmt struct {
	Relations   *List
	RestartSeqs bool
	Behavior    DropBehavior
}

func (n *TruncateStmt) Pos() int {
	return 0
}

func (n *TruncateStmt) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString("TRUNCATE ")
	buf.astFormat(n.Relations)
}
