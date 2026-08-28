package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type CollateExpr struct {
	Tag NodeTag[CollateExpr] `json:"tag"`

	Xpr      Node
	Arg      Node
	CollOid  Oid
	Location int
}

func (n *CollateExpr) Pos() int {
	return n.Location
}

func (n *CollateExpr) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.astFormat(n.Xpr, d)
	buf.WriteString(" COLLATE ")
	buf.astFormat(n.Arg, d)
}
