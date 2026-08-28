package ast

type RangeTableSample struct {
	Tag NodeTag[RangeTableSample] `json:"tag"`

	Relation   Node
	Method     *List
	Args       *List
	Repeatable Node
	Location   int
}

func (n *RangeTableSample) Pos() int {
	return n.Location
}
