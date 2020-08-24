package ast

type SortGroupClause struct {
	TleSortGroupRef Index
	Eqop            Oid
	Sortop          Oid
	NullsFirst      bool
	Hashable        bool
}

func (n *SortGroupClause) Pos() int {
	return 0
}
