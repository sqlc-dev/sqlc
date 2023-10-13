package ast

type Float struct {
	Str string
}

func (n *Float) Pos() int {
	return 0
}

func (n *Float) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString(n.Str)
}
