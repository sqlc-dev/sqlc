package ast

type PartitionBoundSpec struct {
	Tag NodeTag[PartitionBoundSpec] `json:"tag"`

	Strategy    byte
	Listdatums  *List `json:",omitempty"`
	Lowerdatums *List `json:",omitempty"`
	Upperdatums *List `json:",omitempty"`
	Location    int
}

func (n *PartitionBoundSpec) Pos() int {
	return n.Location
}
