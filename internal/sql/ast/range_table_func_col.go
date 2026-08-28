package ast

type RangeTableFuncCol struct {
	Tag NodeTag[RangeTableFuncCol] `json:"tag"`

	Colname       *string   `json:"colname,omitempty"`
	TypeName      *TypeName `json:"type_name,omitempty"`
	ForOrdinality bool      `json:"for_ordinality"`
	IsNotNull     bool      `json:"is_not_null"`
	Colexpr       Node      `json:"colexpr,omitempty"`
	Coldefexpr    Node      `json:"coldefexpr,omitempty"`
	Location      int       `json:"location"`
}

func (n *RangeTableFuncCol) Pos() int {
	return n.Location
}
