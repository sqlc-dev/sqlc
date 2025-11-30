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
	// Format the arg and type to strings first
	argBuf := NewTrackedBuffer(buf.formatter)
	argBuf.astFormat(n.Arg)

	typeBuf := NewTrackedBuffer(buf.formatter)
	typeBuf.astFormat(n.TypeName)

	buf.WriteString(buf.Cast(argBuf.String(), typeBuf.String()))
}
