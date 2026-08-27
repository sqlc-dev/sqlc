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
	// Beyond the libpg_query set: joins SQLite spells (and treats)
	// distinctly. Both behave as inner joins, but CROSS JOIN carries a
	// planner hint — SQLite will not reorder the pair — and a
	// comma-separated FROM item is its own syntax, so neither may be
	// rewritten into the other.
	JoinTypeCross
	JoinTypeComma
)

type JoinType uint

func (n *JoinType) Pos() int {
	return 0
}
