package ast

type BooleanTest struct {
	Tag NodeTag[BooleanTest] `json:"tag"`

	Xpr          Node `json:",omitempty"`
	Arg          Node `json:",omitempty"`
	Booltesttype BoolTestType
	Location     int
}

func (n *BooleanTest) Pos() int {
	return n.Location
}
