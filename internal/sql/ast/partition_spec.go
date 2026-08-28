package ast

type PartitionSpec struct {
	Tag NodeTag[PartitionSpec] `json:"tag"`

	Strategy   *string `json:"strategy,omitempty"`
	PartParams *List   `json:"part_params,omitempty"`
	Location   int     `json:"location"`
}

func (n *PartitionSpec) Pos() int {
	return n.Location
}
