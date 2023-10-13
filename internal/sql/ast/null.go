package ast

type Null struct {
}

func (n *Null) Pos() int {
	return 0
}
func (n *Null) Format(buf *TrackedBuffer) {
	buf.WriteString("NULL")
}
