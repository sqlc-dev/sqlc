package ast

type PartitionSpec struct {
	Strategy   *string
	PartParams *List
	Location   int
}

func (n *PartitionSpec) Pos() int {
	return n.Location
}
