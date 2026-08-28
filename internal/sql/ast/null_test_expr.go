package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type NullTest struct {
	Tag NodeTag[NullTest] `json:"tag"`

	Xpr          Node         `json:"xpr,omitempty"`
	Arg          Node         `json:"arg,omitempty"`
	Nulltesttype NullTestType `json:"nulltesttype"`
	Argisrow     bool         `json:"argisrow"`
	Location     int          `json:"location"`
}

func (n *NullTest) Pos() int {
	return n.Location
}

// NullTestType values
const (
	NullTestTypeIsNull    NullTestType = 1
	NullTestTypeIsNotNull NullTestType = 2
)

func (n *NullTest) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.astFormat(n.Arg, d)
	switch n.Nulltesttype {
	case NullTestTypeIsNull:
		buf.WriteString(" IS NULL")
	case NullTestTypeIsNotNull:
		buf.WriteString(" IS NOT NULL")
	}
}
