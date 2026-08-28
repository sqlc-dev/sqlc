package ast

type BooleanTest struct {
	Tag NodeTag[BooleanTest] `json:"tag"`

	Xpr          Node
	Arg          Node
	Booltesttype BoolTestType
	Location     int
}

func (n *BooleanTest) Pos() int {
	return n.Location
}
