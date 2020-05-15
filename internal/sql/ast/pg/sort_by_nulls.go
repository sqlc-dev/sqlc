package pg

type SortByNulls uint

func (n *SortByNulls) Pos() int {
	return 0
}
