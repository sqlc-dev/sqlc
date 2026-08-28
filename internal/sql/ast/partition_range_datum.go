package ast

type PartitionRangeDatum struct {
	Tag NodeTag[PartitionRangeDatum] `json:"tag"`

	Kind     PartitionRangeDatumKind `json:"kind"`
	Value    Node                    `json:"value,omitempty"`
	Location int                     `json:"location"`
}

func (n *PartitionRangeDatum) Pos() int {
	return n.Location
}
