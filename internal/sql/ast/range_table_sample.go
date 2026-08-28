package ast

type RangeTableSample struct {
	Tag NodeTag[RangeTableSample] `json:"tag"`

	Relation   Node  `json:",omitempty"`
	Method     *List `json:",omitempty"`
	Args       *List `json:",omitempty"`
	Repeatable Node  `json:",omitempty"`
	Location   int
}

func (n *RangeTableSample) Pos() int {
	return n.Location
}
