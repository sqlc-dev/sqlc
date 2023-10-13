package ast

type SortByDir uint

func (n *SortByDir) Pos() int {
	return 0
}

const (
	SortByDirUndefined SortByDir = 0
	SortByDirDefault   SortByDir = 1
	SortByDirAsc       SortByDir = 2
	SortByDirDesc      SortByDir = 3
	SortByDirUsing     SortByDir = 4
)
