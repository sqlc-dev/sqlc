package ast

import ()

type PartitionRangeDatum struct {
	Kind     PartitionRangeDatumKind
	Value    Node
	Location int
}

func (n *PartitionRangeDatum) Pos() int {
	return n.Location
}
