package pg

type NullTestType uint

func (n *NullTestType) Pos() int {
	return 0
}
