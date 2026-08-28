package ast

type PartitionElem struct {
	Tag NodeTag[PartitionElem] `json:"tag"`

	Name      *string `json:"name,omitempty"`
	Expr      Node    `json:"expr,omitempty"`
	Collation *List   `json:"collation,omitempty"`
	Opclass   *List   `json:"opclass,omitempty"`
	Location  int     `json:"location"`
}

func (n *PartitionElem) Pos() int {
	return n.Location
}
