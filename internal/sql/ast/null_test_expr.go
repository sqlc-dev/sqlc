package ast

type NullTest struct {
	Xpr          Node
	Arg          Node
	Nulltesttype NullTestType
	Argisrow     bool
	Location     int
}

func (n *NullTest) Pos() int {
	return n.Location
}
