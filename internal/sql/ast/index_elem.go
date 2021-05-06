package ast

type IndexElem struct {
	Name          *string
	Expr          Node
	Indexcolname  *string
	Collation     *List
	Opclass       *List
	Ordering      SortByDir
	NullsOrdering SortByNulls
}

func (n *IndexElem) Pos() int {
	return 0
}
