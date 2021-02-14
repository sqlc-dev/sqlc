package ast

type BooleanTest struct {
	Xpr          Node
	Arg          Node
	Booltesttype BoolTestType
	Location     int
}

func (n *BooleanTest) Pos() int {
	return n.Location
}
