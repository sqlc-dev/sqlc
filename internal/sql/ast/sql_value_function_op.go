package ast

type SQLValueFunctionOp uint

const (
	// https://github.com/pganalyze/libpg_query/blob/15-latest/protobuf/pg_query.proto#L2984C1-L3003C1
	_ SQLValueFunctionOp = iota
	SVFOpCurrentDate
	SVFOpCurrentTime
	SVFOpCurrentTimeN
	SVFOpCurrentTimestamp
	SVFOpCurrentTimestampN
	SVFOpLocaltime
	SVFOpLocaltimeN
	SVFOpLocaltimestamp
	SVFOpLocaltimestampN
	SVFOpCurrentRole
	SVFOpCurrentUser
	SVFOpUser
	SVFOpSessionUser
	SVFOpCurrentCatalog
	SVFOpCurrentSchema
)

func (n *SQLValueFunctionOp) Pos() int {
	return 0
}
