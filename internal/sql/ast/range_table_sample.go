package ast

type RangeTableSample struct {
	Tag NodeTag[RangeTableSample] `json:"tag"`

	Relation   Node  `json:"relation,omitempty"`
	Method     *List `json:"method,omitempty"`
	Args       *List `json:"args,omitempty"`
	Repeatable Node  `json:"repeatable,omitempty"`
	Location   int   `json:"location"`
}

func (n *RangeTableSample) Pos() int {
	return n.Location
}
