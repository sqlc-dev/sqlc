package ast

type ArrayRef struct {
	Tag NodeTag[ArrayRef] `json:"tag"`

	Xpr             Node  `json:"xpr,omitempty"`
	Refarraytype    Oid   `json:"refarraytype"`
	Refelemtype     Oid   `json:"refelemtype"`
	Reftypmod       int32 `json:"reftypmod"`
	Refcollid       Oid   `json:"refcollid"`
	Refupperindexpr *List `json:"refupperindexpr,omitempty"`
	Reflowerindexpr *List `json:"reflowerindexpr,omitempty"`
	Refexpr         Node  `json:"refexpr,omitempty"`
	Refassgnexpr    Node  `json:"refassgnexpr,omitempty"`
}

func (n *ArrayRef) Pos() int {
	return 0
}
