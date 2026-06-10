package ast

// DropBehavior captures whether a Postgres DROP statement was qualified with
// RESTRICT or CASCADE. The numeric values mirror pganalyze/pg_query_go's
// DropBehavior enum:
//
//	DropBehavior_UNDEFINED = 0  (no behavior word supplied — same as RESTRICT in PG)
//	DROP_RESTRICT          = 1
//	DROP_CASCADE           = 2
type DropBehavior uint

const (
	DropBehaviorUndefined DropBehavior = 0
	DropBehaviorRestrict  DropBehavior = 1
	DropBehaviorCascade   DropBehavior = 2
)

func (n *DropBehavior) Pos() int {
	return 0
}
