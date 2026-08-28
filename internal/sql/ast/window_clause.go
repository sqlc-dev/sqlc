package ast

type WindowClause struct {
	Tag NodeTag[WindowClause] `json:"tag"`

	Name            *string `json:"name,omitempty"`
	Refname         *string `json:"refname,omitempty"`
	PartitionClause *List   `json:"partition_clause,omitempty"`
	OrderClause     *List   `json:"order_clause,omitempty"`
	FrameOptions    int     `json:"frame_options"`
	StartOffset     Node    `json:"start_offset,omitempty"`
	EndOffset       Node    `json:"end_offset,omitempty"`
	Winref          Index   `json:"winref"`
	CopiedOrder     bool    `json:"copied_order"`
}

func (n *WindowClause) Pos() int {
	return 0
}
