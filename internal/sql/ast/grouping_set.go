package ast

type GroupingSet struct {
	Kind     GroupingSetKind
	Content  *List
	Location int
}

func (n *GroupingSet) Pos() int {
	return n.Location
}
