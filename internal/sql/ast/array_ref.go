package ast

type ArrayRef struct {
	Xpr             Node
	Refarraytype    Oid
	Refelemtype     Oid
	Reftypmod       int32
	Refcollid       Oid
	Refupperindexpr *List
	Reflowerindexpr *List
	Refexpr         Node
	Refassgnexpr    Node
}

func (n *ArrayRef) Pos() int {
	return 0
}
