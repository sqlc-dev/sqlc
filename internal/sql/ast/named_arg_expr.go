package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type NamedArgExpr struct {
	Tag NodeTag[NamedArgExpr] `json:"tag"`

	Xpr       Node    `json:"xpr,omitempty"`
	Arg       Node    `json:"arg,omitempty"`
	Name      *string `json:"name,omitempty"`
	Argnumber int     `json:"argnumber"`
	Location  int     `json:"location"`
}

func (n *NamedArgExpr) Pos() int {
	return n.Location
}

func (n *NamedArgExpr) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	if n.Name != nil {
		buf.WriteString(*n.Name)
	}
	buf.WriteString(" => ")
	buf.astFormat(n.Arg, d)
}
