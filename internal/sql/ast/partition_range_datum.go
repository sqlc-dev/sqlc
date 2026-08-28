package ast

type PartitionRangeDatum struct {
	Tag NodeTag[PartitionRangeDatum] `json:"tag"`

	Kind     PartitionRangeDatumKind
	Value    Node
	Location int
}

func (n *PartitionRangeDatum) Pos() int {
	return n.Location
}
