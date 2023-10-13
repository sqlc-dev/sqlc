package ast

type SortByNulls uint

func (n *SortByNulls) Pos() int {
	return 0
}

const (
	SortByNullsUndefined SortByNulls = 0
	SortByNullsDefault   SortByNulls = 1
	SortByNullsFirst     SortByNulls = 2
	SortByNullsLast      SortByNulls = 3
)
