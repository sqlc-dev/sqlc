package ast

type IntoClause struct {
	Tag NodeTag[IntoClause] `json:"tag"`

	Rel            *RangeVar      `json:"rel,omitempty"`
	ColNames       *List          `json:"col_names,omitempty"`
	Options        *List          `json:"options,omitempty"`
	OnCommit       OnCommitAction `json:"on_commit"`
	TableSpaceName *string        `json:"table_space_name,omitempty"`
	ViewQuery      Node           `json:"view_query,omitempty"`
	SkipData       bool           `json:"skip_data"`
}

func (n *IntoClause) Pos() int {
	return 0
}
