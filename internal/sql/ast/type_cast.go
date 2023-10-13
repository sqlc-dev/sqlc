package ast

type TypeCast struct {
	Arg      Node
	TypeName *TypeName
	Location int
}

func (n *TypeCast) Pos() int {
	return n.Location
}

func (n *TypeCast) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.astFormat(n.Arg)
	buf.WriteString("::")
	buf.astFormat(n.TypeName)
}
