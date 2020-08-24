package ast

type NullTestType uint

func (n *NullTestType) Pos() int {
	return 0
}
