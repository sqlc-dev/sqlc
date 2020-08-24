package ast

type PartitionCmd struct {
	Name  *RangeVar
	Bound *PartitionBoundSpec
}

func (n *PartitionCmd) Pos() int {
	return 0
}
