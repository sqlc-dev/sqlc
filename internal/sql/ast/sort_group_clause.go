package ast

type SortGroupClause struct {
	Tag NodeTag[SortGroupClause] `json:"tag"`

	TleSortGroupRef Index `json:"tle_sort_group_ref"`
	Eqop            Oid   `json:"eqop"`
	Sortop          Oid   `json:"sortop"`
	NullsFirst      bool  `json:"nulls_first"`
	Hashable        bool  `json:"hashable"`
}

func (n *SortGroupClause) Pos() int {
	return 0
}
