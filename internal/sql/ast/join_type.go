package ast

// JoinType is the reported type of the join
// Enum copies https://github.com/pganalyze/libpg_query/blob/13-latest/protobuf/pg_query.proto#L2890-L2901
const (
	_ JoinType = iota
	JoinTypeInner
	JoinTypeLeft
	JoinTypeFull
	JoinTypeRight
	JoinTypeSemi
	JoinTypeAnti
	JoinTypeUniqueOuter
	JoinTypeUniqueInner
)

type JoinType uint

func (n *JoinType) Pos() int {
	return 0
}
