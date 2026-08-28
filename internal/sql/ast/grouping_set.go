package ast

type GroupingSet struct {
	Tag NodeTag[GroupingSet] `json:"tag"`

	Kind     GroupingSetKind `json:"kind"`
	Content  *List           `json:"content,omitempty"`
	Location int             `json:"location"`
}

func (n *GroupingSet) Pos() int {
	return n.Location
}
