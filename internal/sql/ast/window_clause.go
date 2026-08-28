package ast

type WindowClause struct {
	Tag NodeTag[WindowClause] `json:"tag"`

	Name            *string `json:",omitempty"`
	Refname         *string `json:",omitempty"`
	PartitionClause *List   `json:",omitempty"`
	OrderClause     *List   `json:",omitempty"`
	FrameOptions    int
	StartOffset     Node `json:",omitempty"`
	EndOffset       Node `json:",omitempty"`
	Winref          Index
	CopiedOrder     bool
}

func (n *WindowClause) Pos() int {
	return 0
}
