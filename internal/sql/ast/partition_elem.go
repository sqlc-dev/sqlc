package ast

type PartitionElem struct {
	Tag NodeTag[PartitionElem] `json:"tag"`

	Name      *string `json:",omitempty"`
	Expr      Node    `json:",omitempty"`
	Collation *List   `json:",omitempty"`
	Opclass   *List   `json:",omitempty"`
	Location  int
}

func (n *PartitionElem) Pos() int {
	return n.Location
}
