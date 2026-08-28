package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type SQLValueFunction struct {
	Tag NodeTag[SQLValueFunction] `json:"tag"`

	Xpr      Node               `json:"xpr,omitempty"`
	Op       SQLValueFunctionOp `json:"op"`
	Type     Oid                `json:"type"`
	Typmod   int32              `json:"typmod"`
	Location int                `json:"location"`
}

func (n *SQLValueFunction) Pos() int {
	return n.Location
}

func (n *SQLValueFunction) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	switch n.Op {
	case SVFOpCurrentDate:
		buf.WriteString("CURRENT_DATE")
	case SVFOpCurrentTime:
	case SVFOpCurrentTimeN:
	case SVFOpCurrentTimestamp:
	case SVFOpCurrentTimestampN:
	case SVFOpLocaltime:
	case SVFOpLocaltimeN:
	case SVFOpLocaltimestamp:
	case SVFOpLocaltimestampN:
	case SVFOpCurrentRole:
	case SVFOpCurrentUser:
	case SVFOpUser:
	case SVFOpSessionUser:
	case SVFOpCurrentCatalog:
	case SVFOpCurrentSchema:
	}
}
