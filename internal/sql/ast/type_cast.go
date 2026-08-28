package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type TypeCast struct {
	Tag NodeTag[TypeCast] `json:"tag"`

	Arg      Node      `json:"arg,omitempty"`
	TypeName *TypeName `json:"type_name,omitempty"`
	Location int       `json:"location"`
}

func (n *TypeCast) Pos() int {
	return n.Location
}

func (n *TypeCast) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	// Format the arg and type to strings first
	argBuf := NewTrackedBuffer()
	argBuf.astFormat(n.Arg, d)

	typeBuf := NewTrackedBuffer()
	typeBuf.astFormat(n.TypeName, d)

	buf.WriteString(d.Cast(argBuf.String(), typeBuf.String()))
}
