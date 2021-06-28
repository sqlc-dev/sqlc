package ast

type InferClause struct {
	IndexElems  *List
	WhereClause Node
	Conname     *string
	Location    int
}

func (n *InferClause) Pos() int {
	return n.Location
}
