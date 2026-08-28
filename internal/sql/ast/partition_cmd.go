package ast

type PartitionCmd struct {
	Tag NodeTag[PartitionCmd] `json:"tag"`

	Name  *RangeVar           `json:"name,omitempty"`
	Bound *PartitionBoundSpec `json:"bound,omitempty"`
}

func (n *PartitionCmd) Pos() int {
	return 0
}
