package ast

type IntoClause struct {
	Tag NodeTag[IntoClause] `json:"tag"`

	Rel            *RangeVar `json:",omitempty"`
	ColNames       *List     `json:",omitempty"`
	Options        *List     `json:",omitempty"`
	OnCommit       OnCommitAction
	TableSpaceName *string `json:",omitempty"`
	ViewQuery      Node    `json:",omitempty"`
	SkipData       bool
}

func (n *IntoClause) Pos() int {
	return 0
}
