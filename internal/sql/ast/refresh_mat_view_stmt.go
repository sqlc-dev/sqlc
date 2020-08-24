package ast

type RefreshMatViewStmt struct {
	Concurrent bool
	SkipData   bool
	Relation   *RangeVar
}

func (n *RefreshMatViewStmt) Pos() int {
	return 0
}
