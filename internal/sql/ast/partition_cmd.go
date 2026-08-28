package ast

type PartitionCmd struct {
	Tag NodeTag[PartitionCmd] `json:"tag"`

	Name  *RangeVar           `json:",omitempty"`
	Bound *PartitionBoundSpec `json:",omitempty"`
}

func (n *PartitionCmd) Pos() int {
	return 0
}
