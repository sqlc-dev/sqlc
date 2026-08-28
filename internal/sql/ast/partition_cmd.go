package ast

type PartitionCmd struct {
	Tag NodeTag[PartitionCmd] `json:"tag"`

	Name  *RangeVar
	Bound *PartitionBoundSpec
}

func (n *PartitionCmd) Pos() int {
	return 0
}
