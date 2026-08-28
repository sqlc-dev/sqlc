package ast

type TableSampleClause struct {
	Tag NodeTag[TableSampleClause] `json:"tag"`

	Tsmhandler Oid   `json:"tsmhandler"`
	Args       *List `json:"args,omitempty"`
	Repeatable Node  `json:"repeatable,omitempty"`
}

func (n *TableSampleClause) Pos() int {
	return 0
}
