package ast

type GroupingSet struct {
	Tag NodeTag[GroupingSet] `json:"tag"`

	Kind     GroupingSetKind
	Content  *List `json:",omitempty"`
	Location int
}

func (n *GroupingSet) Pos() int {
	return n.Location
}
