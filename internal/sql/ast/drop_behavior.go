package ast

type DropBehavior uint

// Matches pganalyze/pg_query_go DropBehavior enum:
// DropBehavior_UNDEFINED = 0, DROP_RESTRICT = 1, DROP_CASCADE = 2.
const (
	DropBehaviorUndefined DropBehavior = 0
	DropBehaviorRestrict  DropBehavior = 1
	DropBehaviorCascade   DropBehavior = 2
)

func (n *DropBehavior) Pos() int {
	return 0
}
