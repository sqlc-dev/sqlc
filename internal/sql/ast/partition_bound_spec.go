package ast

type PartitionBoundSpec struct {
	Tag NodeTag[PartitionBoundSpec] `json:"tag"`

	Strategy    byte  `json:"strategy"`
	Listdatums  *List `json:"listdatums,omitempty"`
	Lowerdatums *List `json:"lowerdatums,omitempty"`
	Upperdatums *List `json:"upperdatums,omitempty"`
	Location    int   `json:"location"`
}

func (n *PartitionBoundSpec) Pos() int {
	return n.Location
}
