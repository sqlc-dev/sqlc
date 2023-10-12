package ast

type SortBy struct {
	Node        Node
	SortbyDir   SortByDir
	SortbyNulls SortByNulls
	UseOp       *List
	Location    int
}

func (n *SortBy) Pos() int {
	return n.Location
}

func (n *SortBy) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.astFormat(n.Node)
	switch n.SortbyDir {
	case SortByDirAsc:
		buf.WriteString(" ASC")
	case SortByDirDesc:
		buf.WriteString(" DESC")
	}
}
