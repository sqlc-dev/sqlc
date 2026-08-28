package ast

type TableSampleClause struct {
	Tag NodeTag[TableSampleClause] `json:"tag"`

	Tsmhandler Oid
	Args       *List `json:",omitempty"`
	Repeatable Node  `json:",omitempty"`
}

func (n *TableSampleClause) Pos() int {
	return 0
}
