package ast

type TableSampleClause struct {
	Tag NodeTag[TableSampleClause] `json:"tag"`

	Tsmhandler Oid
	Args       *List
	Repeatable Node
}

func (n *TableSampleClause) Pos() int {
	return 0
}
