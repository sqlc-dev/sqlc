package ast

type PartitionElem struct {
	Tag NodeTag[PartitionElem] `json:"tag"`

	Name      *string
	Expr      Node
	Collation *List
	Opclass   *List
	Location  int
}

func (n *PartitionElem) Pos() int {
	return n.Location
}
