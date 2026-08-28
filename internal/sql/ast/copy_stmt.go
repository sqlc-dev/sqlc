package ast

type CopyStmt struct {
	Tag NodeTag[CopyStmt] `json:"tag"`

	Relation  *RangeVar
	Query     Node
	Attlist   *List
	IsFrom    bool
	IsProgram bool
	Filename  *string
	Options   *List
}

func (n *CopyStmt) Pos() int {
	return 0
}
