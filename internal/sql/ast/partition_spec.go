package ast

type PartitionSpec struct {
	Tag NodeTag[PartitionSpec] `json:"tag"`

	Strategy   *string `json:",omitempty"`
	PartParams *List   `json:",omitempty"`
	Location   int
}

func (n *PartitionSpec) Pos() int {
	return n.Location
}
