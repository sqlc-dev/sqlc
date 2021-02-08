package ast

type CopyStmt struct {
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
