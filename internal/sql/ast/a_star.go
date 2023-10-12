package ast

type A_Star struct {
}

func (n *A_Star) Pos() int {
	return 0
}

func (n *A_Star) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteRune('*')
}
