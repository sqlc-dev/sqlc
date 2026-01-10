package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

// YDB specific
type Pragma_stmt struct {
	Name     Node
	Cols     *List
	Equals   bool
	Values   *List
	Location int
}

func (n *Pragma_stmt) Pos() int {
	return n.Location
}

func (n *Pragma_stmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}

	buf.WriteString("PRAGMA ")
	if n.Name != nil {
		buf.astFormat(n.Name, d)
	}
	if n.Cols != nil {
		buf.astFormat(n.Cols, d)
	}

	if n.Equals {
		buf.WriteString(" = ")
	}

	if n.Values != nil {
		if n.Equals {
			buf.astFormat(n.Values, d)
		} else {
			buf.WriteString("(")
			buf.astFormat(n.Values, d)
			buf.WriteString(")")
		}
	}

}
