package ast

type ArrayRef struct {
	Tag NodeTag[ArrayRef] `json:"tag"`

	Xpr             Node `json:",omitempty"`
	Refarraytype    Oid
	Refelemtype     Oid
	Reftypmod       int32
	Refcollid       Oid
	Refupperindexpr *List `json:",omitempty"`
	Reflowerindexpr *List `json:",omitempty"`
	Refexpr         Node  `json:",omitempty"`
	Refassgnexpr    Node  `json:",omitempty"`
}

func (n *ArrayRef) Pos() int {
	return 0
}
