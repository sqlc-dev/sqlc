package ast

type WindowClause struct {
	Tag NodeTag[WindowClause] `json:"tag"`

	Name            *string
	Refname         *string
	PartitionClause *List
	OrderClause     *List
	FrameOptions    int
	StartOffset     Node
	EndOffset       Node
	Winref          Index
	CopiedOrder     bool
}

func (n *WindowClause) Pos() int {
	return 0
}
