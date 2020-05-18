package sqlc

type ColumnRef struct {
	Name string
}

func (n *ColumnRef) Pos() int {
	return 0
}
