package ast

type IntoClause struct {
	Tag NodeTag[IntoClause] `json:"tag"`

	Rel            *RangeVar
	ColNames       *List
	Options        *List
	OnCommit       OnCommitAction
	TableSpaceName *string
	ViewQuery      Node
	SkipData       bool
}

func (n *IntoClause) Pos() int {
	return 0
}
