package ast

type GroupingFunc struct {
	Tag NodeTag[GroupingFunc] `json:"tag"`

	Xpr         Node  `json:",omitempty"`
	Args        *List `json:",omitempty"`
	Refs        *List `json:",omitempty"`
	Cols        *List `json:",omitempty"`
	Agglevelsup Index
	Location    int
}

func (n *GroupingFunc) Pos() int {
	return n.Location
}
