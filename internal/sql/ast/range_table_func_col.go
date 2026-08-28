package ast

type RangeTableFuncCol struct {
	Tag NodeTag[RangeTableFuncCol] `json:"tag"`

	Colname       *string   `json:",omitempty"`
	TypeName      *TypeName `json:",omitempty"`
	ForOrdinality bool
	IsNotNull     bool
	Colexpr       Node `json:",omitempty"`
	Coldefexpr    Node `json:",omitempty"`
	Location      int
}

func (n *RangeTableFuncCol) Pos() int {
	return n.Location
}
