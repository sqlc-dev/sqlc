package ast

type PartitionSpec struct {
	Tag NodeTag[PartitionSpec] `json:"tag"`

	Strategy   *string
	PartParams *List
	Location   int
}

func (n *PartitionSpec) Pos() int {
	return n.Location
}
