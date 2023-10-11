package ast

type SQLValueFunction struct {
	Xpr      Node
	Op       SQLValueFunctionOp
	Type     Oid
	Typmod   int32
	Location int
}

func (n *SQLValueFunction) Pos() int {
	return n.Location
}

func (n *SQLValueFunction) Format(buf *TrackedBuffer) {
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
