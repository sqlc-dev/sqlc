package ast

import "strconv"

type Integer struct {
	Ival int64
}

func (n *Integer) Pos() int {
	return 0
}

func (n *Integer) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString(strconv.FormatInt(n.Ival, 10))
}
