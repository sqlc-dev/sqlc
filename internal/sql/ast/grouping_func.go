package ast

type GroupingFunc struct {
	Tag NodeTag[GroupingFunc] `json:"tag"`

	Xpr         Node  `json:"xpr,omitempty"`
	Args        *List `json:"args,omitempty"`
	Refs        *List `json:"refs,omitempty"`
	Cols        *List `json:"cols,omitempty"`
	Agglevelsup Index `json:"agglevelsup"`
	Location    int   `json:"location"`
}

func (n *GroupingFunc) Pos() int {
	return n.Location
}
