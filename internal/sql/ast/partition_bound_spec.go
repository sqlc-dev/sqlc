package ast

type PartitionBoundSpec struct {
	Tag NodeTag[PartitionBoundSpec] `json:"tag"`

	Strategy    byte
	Listdatums  *List
	Lowerdatums *List
	Upperdatums *List
	Location    int
}

func (n *PartitionBoundSpec) Pos() int {
	return n.Location
}
