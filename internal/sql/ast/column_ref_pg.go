package ast

type ColumnRef_PG struct {
	Fields   *List
	Location int
}

func (n *ColumnRef_PG) Pos() int {
	return n.Location
}
