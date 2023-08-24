package ast

// https://github.com/pganalyze/libpg_query/blob/13-latest/protobuf/pg_query.proto#L2783-L2789
const (
	_ BoolExprType = iota
	BoolExprTypeAnd
	BoolExprTypeOr
	BoolExprTypeNot

	// Added for MySQL
	BoolExprTypeIsNull
	BoolExprTypeIsNotNull
)

type BoolExprType uint

func (n *BoolExprType) Pos() int {
	return 0
}
