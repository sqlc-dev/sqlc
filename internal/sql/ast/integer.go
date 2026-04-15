package ast

import (
	"strconv"

	"github.com/sqlc-dev/sqlc/internal/sql/format"
)

type Integer struct {
	Ival int64
}

func (n *Integer) Pos() int {
	return 0
}

func (n *Integer) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString(strconv.FormatInt(n.Ival, 10))
}
