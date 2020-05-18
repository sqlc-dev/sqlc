package ast

type WindowDef struct {
	Name            *string
	Refname         *string
	PartitionClause *List
	OrderClause     *List
	FrameOptions    int
	StartOffset     Node
	EndOffset       Node
	Location        int
}

func (n *WindowDef) Pos() int {
	return n.Location
}
