package ast

type BooleanTest struct {
	Tag NodeTag[BooleanTest] `json:"tag"`

	Xpr          Node         `json:"xpr,omitempty"`
	Arg          Node         `json:"arg,omitempty"`
	Booltesttype BoolTestType `json:"booltesttype"`
	Location     int          `json:"location"`
}

func (n *BooleanTest) Pos() int {
	return n.Location
}
