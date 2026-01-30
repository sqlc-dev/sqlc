package ast

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/sql/format"
)

type ParamRef struct {
	Number   int
	Location int
	Dollar   bool

	// YDB specific
	Plike bool
}

func (n *ParamRef) Pos() int {
	return n.Location
}

func (n *ParamRef) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	if n.Plike {
		fmt.Fprintf(buf, "$p%d", n.Number)
		return
	}
	buf.WriteString(d.Param(n.Number))
}
